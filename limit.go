package limiter

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type LimitController struct {
	RedisDB *redis.Client
	GIpRate GlobalRate
	IpRate  []singleRate
	Record  bool
}

func DefaultController(rdb *redis.Client, gIp GlobalRate) (*LimitController, error) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis server doesn't collect!")
		return nil, err
	}

	return &LimitController{
		RedisDB: rdb,
		GIpRate: gIp,
		IpRate:  nil,
		Record:  false,
	}, nil
}
