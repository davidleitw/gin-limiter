package limiter

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

type LimitController struct {
	RedisDB     *redis.Client
	globalRate  *GlobalRate
	routerRates Rates
	Record      bool
	script      string
	mode        string
	logger      *log.Logger
}

func createController(rdb *redis.Client, gr *GlobalRate, sr Rates, record bool, mode string) *LimitController {
	var logger *log.Logger
	if mode == "debug" {
		logger = log.New(os.Stdout, "[Limit] ", log.Ldate|log.Ltime)
	}

	return &LimitController{
		RedisDB:     rdb,
		globalRate:  gr,
		routerRates: sr,
		Record:      record,
		script:      "",
		mode:        mode,
		logger:      logger,
	}
}

func DefaultController(rdb *redis.Client, command string, limit int, mode string) (*LimitController, error) {
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

	return createController(rdb, gRate, Rates{}, false, mode), nil
}

func (lc *LimitController) Mode() string {
	return lc.mode
}

func (lc *LimitController) SetShaScript(sha string) {
	lc.script = sha
}

func (lc *LimitController) GetShaScript() string {
	return lc.script
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

// 根據router資訊新增一個對於該router的limiter
func (lc *LimitController) Add(path, method, command string, limit int) error {
	sRate, err := newSingleRate(path, method, command, limit)
	if err != nil {
		return err
	}

	lc.routerRates.Append(sRate)
	return nil
}

func (lc *LimitController) Init() {
	lc.globalRate.UpdateDeadLine()
	lc.routerRates.UpdateAllDeadLine()

	SHA, err := lc.RedisDB.ScriptLoad(context.Background(), Script).Result()
	if err != nil {
		fmt.Println(err)
	}

	lc.SetShaScript(SHA)
}
