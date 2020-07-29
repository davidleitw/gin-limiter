package limiter

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type LimitController struct {
	RedisDB     *redis.Client
	globalRate  GlobalRate
	routerRates Rates
	Record      bool
}

func createController(rdb *redis.Client, gr GlobalRate, sr []singleRate, record bool) *LimitController {
	return &LimitController{
		RedisDB:     rdb,
		globalRate:  gr,
		routerRates: sr,
		Record:      record,
	}
}

func DefaultController(rdb *redis.Client, gr GlobalRate) (*LimitController, error) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis server doesn't collect!")
		return nil, err
	}

	return createController(rdb, gr, nil, false), nil
}

func (lc *LimitController) UpdateGlobalRate(gr GlobalRate) {
	lc.globalRate = gr
}
