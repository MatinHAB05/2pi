package scraper

import (
	"context"
	"errors"
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

type herLifeConfig struct {
	AllowedDomain    string
	BaseURL          string
	BlogStartURL     string
	UserAgent        string
	RedisPrefix      string
	NumWorkers       int
	QueueSize        int
	Parallelism      int
	Delay            time.Duration
	CategorySelector string
	ArticleSelector  string
}

func NewHerLifeConfig() herLifeConfig {
	return herLifeConfig{}
}

func DefaultHerLifeConfig() herLifeConfig {
	return herLifeConfig{
		AllowedDomain:    "herlifeapp.com",
		BaseURL:          "https://herlifeapp.com",
		BlogStartURL:     "https://herlifeapp.com/blog/",
		UserAgent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		RedisPrefix:      "colly_herlife",
		NumWorkers:       5,
		QueueSize:        100,
		Parallelism:      10,
		Delay:            100 * time.Millisecond,
		CategorySelector: `a[href*="/blog/category/"]`,
		ArticleSelector:  `a[href*="/blog/articles/"]`,
	}
}

type herLifeArticle struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
}

type herLifeEngine struct {
	cfg      *herLifeConfig
	redisCfg config.Redis
	logger   logger.Logger
	repo     repository_contract.ArticleRepository
	esrepo   repository_contract.ArticleESRepository
}

func NewHerLifeScrapper(
	cfg *herLifeConfig,
	redisCfg config.Redis,
	logger logger.Logger,
	repo repository_contract.ArticleRepository,
	esrepo repository_contract.ArticleESRepository,
) Scraper {
	defaults := DefaultHerLifeConfig()

	if cfg.AllowedDomain == "" {
		cfg.AllowedDomain = defaults.AllowedDomain
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaults.BaseURL
	}
	if cfg.BlogStartURL == "" {
		cfg.BlogStartURL = defaults.BlogStartURL
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = defaults.UserAgent
	}
	if cfg.RedisPrefix == "" {
		cfg.RedisPrefix = defaults.RedisPrefix
	}
	if cfg.NumWorkers <= 0 {
		cfg.NumWorkers = defaults.NumWorkers
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = defaults.QueueSize
	}
	if cfg.Parallelism <= 0 {
		cfg.Parallelism = defaults.Parallelism
	}
	if cfg.Delay <= 0 {
		cfg.Delay = defaults.Delay
	}
	if cfg.CategorySelector == "" {
		cfg.CategorySelector = defaults.CategorySelector
	}
	if cfg.ArticleSelector == "" {
		cfg.ArticleSelector = defaults.ArticleSelector
	}

	return &herLifeEngine{
		cfg:      cfg,
		redisCfg: redisCfg,
		logger:   logger,
		repo:     repo,
		esrepo:   esrepo,
	}
}

func (engine *herLifeEngine) WhoIsTarget() TargetInfo {
	return TargetInfo{
		Name:    "HerLife-Blogs",
		MainURL: engine.cfg.BlogStartURL,
	}
}

func (engine *herLifeEngine) Run(ctx context.Context, scrapDataMu *sync.Mutex) error {
	return engine.herlife(ctx, scrapDataMu)
}

func (engine *herLifeEngine) herlife(ctx context.Context, scrapDataMu *sync.Mutex) error {
	allowedDomain := engine.cfg.AllowedDomain

	storage := &redisstorage.Storage{
		Address:  engine.redisCfg.Host + ":" + strconv.Itoa(engine.redisCfg.Port),
		Password: engine.redisCfg.Password,
		DB:       engine.redisCfg.RDBNumber,
		Prefix:   engine.cfg.RedisPrefix,
	}

	baseCollector := colly.NewCollector(
		colly.AllowedDomains(allowedDomain, "www."+allowedDomain),
		colly.Async(true),
		colly.UserAgent(engine.cfg.UserAgent),
	)

	if err := baseCollector.SetStorage(storage); err != nil {
		engine.logger.Errorf("failed to set redis storage for colly: %v", err)
		return fmt.Errorf("failed to set redis storage: %w", err)
	}
	defer storage.Client.Close()

	err := baseCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*" + allowedDomain,
		Parallelism: engine.cfg.Parallelism,
		Delay:       engine.cfg.Delay,
	})
	if err != nil {
		engine.logger.Errorf("failed to set limit rule: %v", err)
		return err
	}

	baseCollector.OnRequest(func(r *colly.Request) {})

	baseCollector.OnError(func(r *colly.Response, err error) {
		engine.logger.Errorf("request failed for URL %s: %v", r.Request.URL, err)
	})

	// -------------------------------------------------------------
	// SETUP QUEUE & WORKER POOL FOR DB SAVES
	// -------------------------------------------------------------
	articleQueue := make(chan herLifeArticle, engine.cfg.QueueSize)
	var dbWg sync.WaitGroup

	// Start worker goroutines to process DB writes concurrently from the queue
	for i := 0; i < engine.cfg.NumWorkers; i++ {
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
					engine.logger.Errorf("[Worker %d] failed to save article in db [%s]: %v", workerID, art.URL, err)
				}

				esdbArticle := &repository_contract.ArticleDocument{
					ID:          strconv.FormatInt(dbArticle.ID, 10),
					CreatedAt:   dbArticle.CreatedAt,
					UpdatedAt:   dbArticle.UpdatedAt,
					DeletedAt:   dbArticle.DeletedAt.Time,
					Title:       dbArticle.Title,
					Description: dbArticle.Description,
					ImageURL:    dbArticle.ImageURL,
					URL:         dbArticle.URL,
				}

				if _, err := engine.esrepo.Index(ctx, esdbArticle); err != nil {
					engine.logger.Errorf("[Worker %d] failed to index article in es-db [%s]: %v", workerID, art.URL, err)
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

	categoryCollector.OnHTML(engine.cfg.CategorySelector, func(e *colly.HTMLElement) {
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

	if err := categoryCollector.Visit(engine.cfg.BlogStartURL); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
		engine.logger.Errorf("failed to visit initial category page: %v", err)
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

	articleCollector.OnHTML(engine.cfg.ArticleSelector, func(e *colly.HTMLElement) {
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

	pageCollector.Wait()

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
