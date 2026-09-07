package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/MatinHAB05/2pi/config"
	"github.com/MatinHAB05/2pi/pkg/logger"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
)

// Elastic defines the interface for accessing Elasticsearch clients and configuration.
type Elastic interface {
	GetClient() *elasticsearch.TypedClient
	GetConst() *config.ElasticConst
}

type ElasticSearchDatabase struct {
	client       *elasticsearch.TypedClient
	elasticConst *config.ElasticConst
}

var (
	elasticOnce     sync.Once
	elasticInstance *ElasticSearchDatabase
)

func NewTypedElasticSearchDatabase(elasticConfig *config.ElasticSearchConfig, elasticConst *config.ElasticConst, applogger logger.Logger) Elastic {
	elasticOnce.Do(func() {
		ctx := context.Background()

		loc, err := time.LoadLocation("Asia/Tehran")
		if err != nil {
			loc = time.Local
		}
		timeStamp := time.Now().In(loc).Format("2006-01-02-15-04-05")
		fileName := fmt.Sprintf("%s%s-%s.log", elasticConfig.LogFilePath, timeStamp, uuid.New().String())
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

		client, err := elasticsearch.NewTyped(
			elasticsearch.WithAddresses(fmt.Sprintf("http://%s:%d", elasticConfig.Host, elasticConfig.Port)),
			elasticsearch.WithLogger(&elastictransport.JSONLogger{ //TODO
				Output:             file,
				EnableRequestBody:  true,
				EnableResponseBody: true,
			}),
			// elastictransport.WithLeveledLogger(NewElasticLoggerAdapter(applogger)), // مسخره بازی - این برای newclient یا new هست که اون تایپی نیست
			elasticsearch.WithBasicAuth("elastic", elasticConfig.Password),
		)

		if err != nil {
			log.Fatalf("failed to create elasticsearch client: %v", err)
		}

		// Verify connection
		res, err := client.Ping().Do(ctx)
		if err != nil {
			log.Fatalf("failed to ping elasticsearch: %v", err)
		}
		if !res {
			log.Fatal("elasticsearch ping returned false")
		}

		elasticInstance = &ElasticSearchDatabase{
			client:       client,
			elasticConst: elasticConst,
		}
	})

	return elasticInstance
}

func (e *ElasticSearchDatabase) GetClient() *elasticsearch.TypedClient {
	return e.client
}

func (e *ElasticSearchDatabase) GetConst() *config.ElasticConst {
	return e.elasticConst
}

type ElasticLoggerAdapter struct {
	customLogger logger.Logger
}

func NewElasticLoggerAdapter(l logger.Logger) elastictransport.LeveledLogger {
	return &ElasticLoggerAdapter{customLogger: l}
}

func (a *ElasticLoggerAdapter) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	a.customLogger.Debugf(msg+" %v", keysAndValues...)
}

func (a *ElasticLoggerAdapter) Info(ctx context.Context, msg string, keysAndValues ...any) {
	a.customLogger.Infof(msg+" %v", keysAndValues...)
}

func (a *ElasticLoggerAdapter) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	a.customLogger.Warnf(msg+" %v", keysAndValues...)
}

func (a *ElasticLoggerAdapter) Error(ctx context.Context, msg string, keysAndValues ...any) {
	a.customLogger.Errorf(msg+" %v", keysAndValues...)
}
