package limiter

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const TimeFormat = "2006-01-02 15:04:05"

// self define error
var (
	LimitError   = errors.New("Limit should > 0.")
	CommandError = errors.New("The command of first number should > 0.")
	FormatError  = errors.New("Please check the format with your input.")
	MethodError  = errors.New("Please check the method is one of http method.")
)

var timeDict = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": time.Hour * 24,
}

type Dispatcher struct {
	limit    int
	deadline int64
	period   time.Duration
}

// create a limit dispatcher object with command and limit request number.
func LimitDispatcher(command string, limit int) (*Dispatcher, error) {
	var period time.Duration

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

	// limit should > 0
	if limit <= 0 {
		return nil, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period := time.Duration(unit) * t
	} else {
		return nil, FormatError
	}

	return &Dispatcher{limit: limit, deadline: 0, period: period}, nil
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

// get the deadline with format 2006-01-02 15:04:05
func (dispatch *Dispatcher) GetDeadLineWithString() string {
	return time.Unix(dispatch.deadline, 0).Format(TimeFormat)
}

func (dispatch *Dispatcher) MiddleWare(command string, limit int) gin.HandlerFunc {
	// get the deadline for global limit
	deadline := dispatch.GetDeadLine()

	return func(ctx *gin.Context) {
		now := time.Now().Unix()
		path := ctx.FullPath()
		method := ctx.Request.Method
		clientIp := ctx.ClientIP()

		routeKey := path + method + clientIp // for single route limit in redis.
		staticKey := clientIp                // for global limit search in redis.

		// mean global limit should be reset.
		if now > deadline {
			dispatch.UpdateDeadLine()

			// run all restart lua script
			ctx.Next()
		}
	}
}
