package db

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
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
