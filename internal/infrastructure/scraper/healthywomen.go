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
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
)

type healthyWomenConfig struct {
	AllowedDomain string
	BaseURL       string
	UserAgent     string
	NumWorkers    int
	QueueSize     int
	Parallelism   int
	Delay         time.Duration
	ScrollCount   int
	CategoryURLs  []string
	URLPatterns   []string
}

func NewHealthyWomenConfig() healthyWomenConfig {
	return healthyWomenConfig{}
}

func DefaultHealthyWomenConfig() healthyWomenConfig {
	return healthyWomenConfig{
		AllowedDomain: "healthywomen.org",
		BaseURL:       "https://www.healthywomen.org",
		UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		NumWorkers:    10,
		QueueSize:     100,
		Parallelism:   20,
		Delay:         100 * time.Millisecond,
		ScrollCount:   5,
		CategoryURLs: []string{
			"https://www.healthywomen.org/your-wellness/",
			"https://www.healthywomen.org/your-health/",
			"https://www.healthywomen.org/your-health/prevention--screenings/",
			"https://www.healthywomen.org/your-health/sexual-health/",
			"https://www.healthywomen.org/your-health/pregnancy--postpartum/",
			"https://www.healthywomen.org/your-health/fertility/",
		},
		URLPatterns: []string{
			"/your-wellness/",
			"/your-health/",
			"/your-care/",
		},
	}
}

type healthyWomenArticle struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
}

type healthyWomenEngine struct {
	cfg      *healthyWomenConfig
	redisCfg config.Redis
	logger   logger.Logger
	repo     repository_contract.ArticleRepository
}

func NewHealthyWomenScraper(
	cfg *healthyWomenConfig,
	redisCfg config.Redis,
	logger logger.Logger,
	repo repository_contract.ArticleRepository,
) Scraper {
	// Apply default values for uninitialized fields
	defaults := DefaultHealthyWomenConfig()
	if cfg.AllowedDomain == "" {
		cfg.AllowedDomain = defaults.AllowedDomain
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaults.BaseURL
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = defaults.UserAgent
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
	if cfg.ScrollCount <= 0 {
		cfg.ScrollCount = defaults.ScrollCount
	}
	if len(cfg.CategoryURLs) == 0 {
		cfg.CategoryURLs = defaults.CategoryURLs
	}
	if len(cfg.URLPatterns) == 0 {
		cfg.URLPatterns = defaults.URLPatterns
	}

	return &healthyWomenEngine{
		cfg:      cfg,
		redisCfg: redisCfg,
		logger:   logger,
		repo:     repo,
	}
}

func (engine *healthyWomenEngine) WhoIsTarget() TargetInfo {
	return TargetInfo{
		Name:    "HealthyWomen-Blogs",
		MainURL: engine.cfg.BaseURL,
	}
}

func (engine *healthyWomenEngine) Run(ctx context.Context, scrapDataMu *sync.Mutex) error {
	return engine.healthyWomen(ctx, scrapDataMu)
}

func (engine *healthyWomenEngine) healthyWomen(ctx context.Context, scrapDataMu *sync.Mutex) error {
	allowedDomain := engine.cfg.AllowedDomain

	storage := &redisstorage.Storage{
		Address:  engine.redisCfg.Host + ":" + strconv.Itoa(engine.redisCfg.Port),
		Password: engine.redisCfg.Password,
		DB:       engine.redisCfg.RDBNumber,
		Prefix:   "colly_healthywomen",
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

	baseCollector.OnError(func(r *colly.Response, err error) {
		engine.logger.Errorf("request failed for URL %s: %v", r.Request.URL, err)
	})

	// -------------------------------------------------------------
	// SETUP QUEUE & WORKER POOL FOR DB SAVES
	// -------------------------------------------------------------
	articleQueue := make(chan healthyWomenArticle, engine.cfg.QueueSize)
	var dbWg sync.WaitGroup

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
					engine.logger.Errorf("[Worker %d] failed to save article [%s]: %v", workerID, art.URL, err)
				}
			}
		}(i)
	}

	// -------------------------------------------------------------
	// STEP 1: Collect Article Links using Chromedp (Infinite Scroll)
	// -------------------------------------------------------------
	articleURLs, err := engine.fetchCategoryArticleURLs(ctx, engine.cfg.CategoryURLs)
	if err != nil {
		engine.logger.Errorf("failed during category infinite scroll extraction: %v", err)
		close(articleQueue)
		return err
	}

	engine.logger.Infof("Found %d unique article URLs across %d categories", len(articleURLs), len(engine.cfg.CategoryURLs))

	// -------------------------------------------------------------
	// STEP 2: Extract Content via Colly & Enqueue Immediately
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

		article := healthyWomenArticle{
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

	engine.logger.Infof("Finished scraping and processing all database writes for HealthyWomen.")
	return nil
}

func (engine *healthyWomenEngine) fetchCategoryArticleURLs(ctx context.Context, categoryURLs []string) ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent(engine.cfg.UserAgent),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	chromeCtx, cancelChrome := chromedp.NewContext(allocCtx)
	defer cancelChrome()

	var artMu sync.Mutex
	visitedURLs := make(map[string]bool)
	var articleURLs []string

	for _, catURL := range categoryURLs {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		engine.logger.Infof("Chromedp navigating and scrolling: %s", catURL)

		var renderedHTML string
		err := chromedp.Run(chromeCtx,
			chromedp.Navigate(catURL),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			engine.logger.Errorf("Chromedp navigation failed for %s: %v", catURL, err)
			continue
		}

		for i := 1; i <= engine.cfg.ScrollCount; i++ {
			_ = chromedp.Run(chromeCtx,
				chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil),
				chromedp.Sleep(2*time.Second),
			)
		}

		if err := chromedp.Run(chromeCtx, chromedp.OuterHTML(`html`, &renderedHTML, chromedp.ByQuery)); err != nil {
			engine.logger.Errorf("Chromedp failed to retrieve HTML for %s: %v", catURL, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(renderedHTML))
		if err != nil {
			engine.logger.Errorf("goquery parse error: %v", err)
			continue
		}

		doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}

			if engine.matchesURLPattern(href) {
				if strings.HasPrefix(href, "/") {
					href = engine.cfg.BaseURL + href
				}

				artMu.Lock()
				if !visitedURLs[href] {
					visitedURLs[href] = true
					articleURLs = append(articleURLs, href)
				}
				artMu.Unlock()
			}
		})
	}

	return articleURLs, nil
}

func (engine *healthyWomenEngine) matchesURLPattern(href string) bool {
	for _, pattern := range engine.cfg.URLPatterns {
		if strings.Contains(href, pattern) {
			return true
		}
	}
	return false
}
