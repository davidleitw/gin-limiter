package limiter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// for the deadline time format.
const TimeFormat = "2006-01-02 15:04:05"

// self define error
var (
	LimitError   = errors.New("Limit should > 0.")
	CommandError = errors.New("The command of first number should > 0.")
	FormatError  = errors.New("Please check the format with your input.")
	MethodError  = errors.New("Please check the method is one of http method.")
	ServerError  = errors.New("StatusInternalServerError, please wait a minute.")
)

var timeDict = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": time.Hour * 24,
}

type Dispatcher struct {
	limit       int
	deadline    int64
	shaScript   map[string]string
	period      time.Duration
	redisClient *redis.Client
}

// create a limit dispatcher object with command and limit request number.
func LimitDispatcher(command string, limit int, rdb *redis.Client) (*Dispatcher, error) {

	dispatcher := new(Dispatcher)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	dispatcher.redisClient = rdb

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return nil, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, FormatError
	}
	if unit <= 0 {
		return nil, CommandError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		dispatcher.period = time.Duration(unit) * t
	} else {
		return nil, FormatError
	}

	// limit should > 0
	if limit <= 0 {
		return nil, LimitError
	}
	dispatcher.limit = limit

	resetSHA, err := dispatcher.redisClient.ScriptLoad(context.Background(), ResetScript).Result()
	if err != nil {
		return nil, err
	}

	normalSHA, err := dispatcher.redisClient.ScriptLoad(context.Background(), Script).Result()
	if err != nil {
		return nil, err
	}

	shaScript := make(map[string]string)
	shaScript["reset"] = resetSHA
	shaScript["normal"] = normalSHA
	dispatcher.shaScript = shaScript
	return dispatcher, nil
}

func (dispatch *Dispatcher) ParseCommand(command string) (time.Duration, error) {
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return period, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return period, FormatError
	}
	if unit <= 0 {
		return period, CommandError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		return time.Duration(unit) * t, nil
	} else {
		return period, FormatError
	}
}

// update the deadline
func (dispatch *Dispatcher) UpdateDeadLine() {
	dispatch.deadline = time.Now().Add(dispatch.period).Unix()
}

// get the limit
func (dispathch *Dispatcher) GetLimit() int {
	return dispathch.limit
}

// get the deadline with unix time.
func (dispatch *Dispatcher) GetDeadLine() int64 {
	return dispatch.deadline
}

func (dispatch *Dispatcher) GetSHAScript(index string) string {
	return dispatch.shaScript[index]
}

// get the deadline with format 2006-01-02 15:04:05
func (dispatch *Dispatcher) GetDeadLineWithString() string {
	return time.Unix(dispatch.deadline, 0).Format(TimeFormat)
}

func (dispatch *Dispatcher) MiddleWare(command string, limit int) gin.HandlerFunc {
	t, _ := dispatch.ParseCommand(command)

	return func(ctx *gin.Context) {
		now := time.Now().Unix()
		clientIp := ctx.ClientIP()
		deadline := dispatch.GetDeadLine()
		routeDeadline := time.Now().Add(t).Unix()
		routeKey := ctx.FullPath() + ctx.Request.Method + clientIp // for single route limit in redis.
		staticKey := clientIp                                      // for global limit search in redis.

		routeLimit := limit
		staticLimit := dispatch.limit

		keys := []string{routeKey, staticKey}
		args := []interface{}{routeLimit, staticLimit, routeDeadline, now}

		// mean global limit should be reset.
		if now > deadline {
			dispatch.UpdateDeadLine()
			_, err := dispatch.redisClient.EvalSha(context.Background(), dispatch.GetSHAScript("reset"), keys, routeDeadline).Result()
			if err != nil {
				log.Println("err = ", err)
				ctx.JSON(http.StatusInternalServerError, err)
				ctx.Abort()
			}
			ctx.Header("X-RateLimit-Limit-global", strconv.Itoa(staticLimit))
			ctx.Header("X-RateLimit-Remaining-global", strconv.Itoa(staticLimit-1))
			ctx.Header("X-RateLimit-Reset-global", dispatch.GetDeadLineWithString())
			ctx.Header("X-RateLimit-Limit-route", strconv.Itoa(limit))
			ctx.Header("X-RateLimit-Remaining-route", strconv.Itoa(limit-1))
			ctx.Header("X-RateLimit-Reset-route", time.Unix(routeDeadline, 0).Format(TimeFormat))
			ctx.Next()
		}

		results, err := dispatch.redisClient.EvalSha(context.Background(), dispatch.GetSHAScript("normal"), keys, args).Result()
		if err != nil {
			log.Println("Result error area, error = ", err)
			ctx.JSON(http.StatusInternalServerError, err)
			ctx.Abort()
		}

		result := results.([]interface{})
		routeRemaining := result[0].(int64)
		staticRemaining := result[1].(int64)
		routedeadline := time.Unix(result[2].(int64), 0).Format(TimeFormat)

		if staticRemaining == -1 {
			ctx.JSON(http.StatusTooManyRequests, dispatch.GetDeadLineWithString())
			ctx.Header("X-RateLimit-Reset-global", dispatch.GetDeadLineWithString())
			ctx.Abort()
		}

		if routeRemaining == -1 {
			ctx.JSON(http.StatusTooManyRequests, routedeadline)
			ctx.Header("X-RateLimit-Reset-single", routedeadline)
			ctx.Abort()
		}

		ctx.Header("X-RateLimit-Limit-global", strconv.Itoa(staticLimit))
		ctx.Header("X-RateLimit-Remaining-global", strconv.FormatInt(staticRemaining, 10))
		ctx.Header("X-RateLimit-Reset-global", dispatch.GetDeadLineWithString())
		ctx.Header("X-RateLimit-Limit-route", strconv.Itoa(routeLimit))
		ctx.Header("X-RateLimit-Remaining-route", strconv.FormatInt(routeRemaining, 10))
		ctx.Header("X-RateLimit-Reset-route", routedeadline)
		ctx.Next()
	}
}
