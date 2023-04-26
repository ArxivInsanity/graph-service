package db

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
	"time"
)

func GetRedisClient() *persist.RedisStore {
	return persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     viper.GetString("redis.connectionUri"),
		Username: viper.GetString("redis.username"),
		Password: os.Getenv("REDIS_CRED"),
		DB:       viper.GetInt("redis.database"),
	}))
}

func GetCachingMiddleware(redisStore *persist.RedisStore) gin.HandlerFunc {
	return cache.CacheByRequestURI(redisStore, 24*time.Hour)
}
