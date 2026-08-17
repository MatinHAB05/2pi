package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/MatinHAB05/2pi/config"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetRDB() *redis.Client
	GetRedisConstant() *config.RedisConst
}

type RedisDatabase struct {
	rdb       *redis.Client
	rdb_const *config.RedisConst
}

var (
	rdbOnce     sync.Once
	rdbInstance *RedisDatabase
)

func NewRedisDatabase(redisConfig *config.Redis, redisConst *config.RedisConst) *RedisDatabase {
	rdbOnce.Do(func() {
		address := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: redisConfig.Password,
			DB:       redisConfig.RDBNumber,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal("Error connecting to Redis:", err)
		}
		rdbInstance = &RedisDatabase{rdb: rdb, rdb_const: redisConst}
	})

	return rdbInstance
}

func (rdb *RedisDatabase) GetRDB() *redis.Client {
	return rdbInstance.rdb
}

func (rdb *RedisDatabase) GetRedisConstant() *config.RedisConst {
	return rdb.rdb_const
}
