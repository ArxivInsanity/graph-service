package db

//
//import (
//	"context"
//	"errors"
//	"github.com/go-redis/redis"
//)
//
//type Redis struct {
//	Client *redis.Client
//}
//
//var (
//	ErrNil = errors.New("no matching record found in redis database")
//	Ctx    = context.
//)
//
//func GetRedisClient(address string) (*Redis, error) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     address,
//		Password: "",
//		DB:       0,
//	})
//	if err := client.Ping(Ctx).Err(); err != nil {
//		return nil, err
//	}
//	return &Redis{
//		Client: client,
//	}, nil
//}
