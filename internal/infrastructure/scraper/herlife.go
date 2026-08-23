package scraper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MatinHAB05/2pi/config"
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	repository_contract "github.com/MatinHAB05/2pi/internal/domain/repository"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
)

type herLifeArticle struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
}

type herLifeEngine struct {
	redisCfg config.Redis
	logger   logger.Logger
	repo     repository_contract.ArticleRepository
}

func NewHerLifeScrapper(redisCfg config.Redis, logger logger.Logger, repo repository_contract.ArticleRepository) Scraper {
	return &herLifeEngine{
		redisCfg: redisCfg,
		logger:   logger,
		repo:     repo,
	}
}

func (engine *herLifeEngine) WhoIsTarget() TargetInfo {
	return TargetInfo{
		Name:    "HerLife-Blogs",
		MainURL: "https://herlifeapp.com/blog",
	}
}

func (engine *herLifeEngine) Run(ctx context.Context, scrapDataMu *sync.Mutex) error {
	return engine.herlife(ctx, scrapDataMu)
}

func (engine *herLifeEngine) herlife(ctx context.Context, scrapDataMu *sync.Mutex) error {
	allowedDomain := "herlifeapp.com"

	storage := &redisstorage.Storage{
		Address:  engine.redisCfg.Host + ":" + strconv.Itoa(engine.redisCfg.Port),
		Password: engine.redisCfg.Password,
		DB:       engine.redisCfg.RDBNumber,
		Prefix:   "colly_herlife",
	}

	baseCollector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain, "www."+allowedDomain),
		colly.Async(true),
	)

	if err := baseCollector.SetStorage(storage); err != nil {
		engine.logger.Errorf("failed to set redis storage for colly:%w", err)
		return fmt.Errorf("failed to set redis storage: %w", err)
	}
	defer storage.Client.Close()

	err := baseCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*" + allowedDomain,
		Parallelism: 10,
		Delay:       100 * time.Millisecond,
	})
	if err != nil {
		engine.logger.Errorf("failed to set limit rule:%w", err)
		return err
	}

	baseCollector.OnRequest(func(r *colly.Request) {})

	baseCollector.OnError(func(r *colly.Response, err error) {
		engine.logger.Errorf("request failed for URL %s:%w", r.Request.URL, err)
	})

	// -------------------------------------------------------------
	// SETUP QUEUE & WORKER POOL FOR DB SAVES
	// -------------------------------------------------------------
	const numWorkers = 5
	const queueSize = 100

	articleQueue := make(chan herLifeArticle, queueSize)
	var dbWg sync.WaitGroup

	// Start worker goroutines to process DB writes concurrently from the queue
	for i := 0; i < numWorkers; i++ {
		dbWg.Add(1)
		go func(workerID int) {
			defer dbWg.Done()
			for art := range articleQueue {
				dbArticle := &entity.Article{
					Title:       art.Title,
					Description: art.Description,
					ImageURL:    art.ImageURL,
					URL:         art.URL,
				}

				if err := engine.repo.Create(ctx, dbArticle); err != nil {
					engine.logger.Errorf("[Worker %d] failed to save article [%s]: %v", workerID, art.URL, err)
				}
			}
		}(i)
	}

	// -------------------------------------------------------------
	// STEP 1: Discover Categories
	// -------------------------------------------------------------
	categoryCollector := baseCollector.Clone()
	var categoryURLs []string
	var catMu sync.Mutex

	categoryCollector.OnHTML(`a[href*="/blog/category/"]`, func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		catMu.Lock()
		defer catMu.Unlock()
		if !engine.contains(categoryURLs, link) {
			categoryURLs = append(categoryURLs, link)
		}
	})

	if err := categoryCollector.Visit("https://herlifeapp.com/blog/"); err != nil {
		engine.logger.Errorf("failed to visit initial category page:%w", err)
		close(articleQueue)
		return err
	}
	categoryCollector.Wait()

	// -------------------------------------------------------------
	// STEP 2: Collect Article Links from Categories
	// -------------------------------------------------------------
	articleCollector := baseCollector.Clone()
	var articleURLs []string
	var artMu sync.Mutex

	articleCollector.OnHTML(`a[href*="/blog/articles/"]`, func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		artMu.Lock()
		defer artMu.Unlock()
		if !engine.contains(articleURLs, link) {
			articleURLs = append(articleURLs, link)
		}
	})

	for _, catURL := range categoryURLs {
		if err := ctx.Err(); err != nil {
			close(articleQueue)
			return err
		}
		_ = articleCollector.Visit(catURL)
	}
	articleCollector.Wait()

	// -------------------------------------------------------------
	// STEP 3: Extract Content & Enqueue Immediately
	// -------------------------------------------------------------
	pageCollector := baseCollector.Clone()

	pageCollector.OnHTML("html", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("h1"))
		if title == "" {
			return
		}

		description := e.ChildAttr(`meta[property="og:description"]`, "content")
		if description == "" {
			description = e.ChildAttr(`meta[name="description"]`, "content")
		}

		imageURL := e.ChildAttr(`meta[property="og:image"]`, "content")

		article := herLifeArticle{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: strings.TrimSpace(description),
			ImageURL:    strings.TrimSpace(imageURL),
		}

		// Push directly to queue as soon as parsed
		select {
		case <-ctx.Done():
			return
		case articleQueue <- article:
		}
	})

	for _, artURL := range articleURLs {
		if err := ctx.Err(); err != nil {
			close(articleQueue)
			return err
		}
		_ = pageCollector.Visit(artURL)
	}

	// Wait for scrapers to complete before closing queue
	pageCollector.Wait()

	// Close channel and wait for DB workers to finish remaining items
	close(articleQueue)
	dbWg.Wait()

	engine.logger.Infof("Finished scraping and processing all database writes.")
	return nil
}

func (engine *herLifeEngine) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
