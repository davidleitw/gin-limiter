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

func createController(rdb *redis.Client, gr GlobalRate, sr Rates, record bool) *LimitController {
	return &LimitController{
		RedisDB:     rdb,
		globalRate:  gr,
		routerRates: sr,
		Record:      record,
	}
}

func DefaultController(rdb *redis.Client, command string, limit int) (*LimitController, error) {
	// Check redis status.
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis server doesn't collect!")
		return nil, err
	}

	// Create GlobalRate object for controller.
	gRate, err := newGlobalRate(command, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Set base limit in redis.
	err = rdb.Set(context.Background(), "limit", gRate.Limit, 0).Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return createController(rdb, gRate, nil, false), nil
}

// lc.UpdateGlobalRate("24-H", 200) => each 24 hours single ip adress can request 200 times for all server router.
func (lc *LimitController) UpdateGlobalRate(command string, limit int) error {
	gRate, err := newGlobalRate(command, limit)
	if err != nil {
		return err
	}

	lc.globalRate = gRate
	return nil
}

func (lc *LimitController) GetGlobalLimit() int {
	return lc.globalRate.Limit
}

// Get single router limit with path.
func (lc *LimitController) GetSingleLimit(path, method string) int {
	return lc.routerRates.getLimit(path, method)
}

func (lc *LimitController) Add(path, command, method string, limit int) error {
	sRate, err := newSingleRate(path, command, method, limit)
	if err != nil {
		return err
	}

	lc.routerRates.Append(sRate)
	return nil
}
