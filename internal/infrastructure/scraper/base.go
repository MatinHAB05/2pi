package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/MatinHAB05/2pi/pkg/logger"
)

// TODO
type Scraper interface {
	WhoIsTarget() TargetInfo
	Run(ctx context.Context, scrapDataMu *sync.Mutex) error
}

type Scrapers struct {
	scrappes    []Scraper
	scrapDataMu *sync.Mutex
}

func NewScrppers(ss []Scraper) Scrapers {
	return Scrapers{
		scrappes:    ss,
		scrapDataMu: &sync.Mutex{},
	}
}

type Article struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
}

type TargetInfo struct {
	Name    string `json:"name"`
	MainURL string `json:"main_url"`
}

// TODO : replace it with database - but for now :
var ScrapData []Article

func (ss Scrapers) RunAll(ctx context.Context, applog logger.Logger) error {

	wg := sync.WaitGroup{}

	for _, s := range ss.scrappes {
		wg.Add(1)
		// Explicitly pass 's' to closure to avoid loop variable capture issues
		go func(sc Scraper) {
			defer wg.Done()

			target := sc.WhoIsTarget()

			// Execute the scraper
			if err := sc.Run(ctx, ss.scrapDataMu); err != nil {
				applog.Error(logger.IO, logger.Scrap, "scraper execution failed", map[logger.ExtraKey]interface{}{
					"target":            target.Name,
					logger.ErrorMessage: err.Error(),
				})
				return
			}

			applog.Info(logger.IO, logger.Scrap, fmt.Sprintf("successfully completed scraping for target: %s", target.Name), nil)
		}(s)
	}

	wg.Wait()

	// -------------------------------------------------------------
	// Save Scraped Data to File
	// -------------------------------------------------------------
	data, err := json.MarshalIndent(ScrapData, "", "  ")
	if err != nil {
		applog.Error(logger.IO, logger.Scrap, "failed to marshal scraped data to JSON", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return fmt.Errorf("marshal indent err: %w", err)
	}

	err = os.WriteFile("for-now.json", data, 0644)
	if err != nil {
		applog.Error(logger.IO, logger.Scrap, "failed to write content to file", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
