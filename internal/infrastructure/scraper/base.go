package scraper

import (
	"context"
	"fmt"
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

	return nil
}
