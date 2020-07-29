package limiter

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// global ip rate
type GlobalRate struct {
	Command string
	Period  time.Duration
	Limit   int
}

// local ip rate
type singleRate struct {
	Path    string
	Command string
	Period  time.Duration
	Limit   int
}

type Rates []singleRate

var timeDict = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": time.Hour * 24,
}

var FormatError = errors.New("please check the format with your input.")
var LimitError = errors.New("limit should > 0.")

// NewGlobalRate("10-m", 200)
// Each 10 minutes single ip address can request 200 times.
func NewGlobalRate(command string, limit int) (GlobalRate, error) {
	var grate GlobalRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!")
		return grate, FormatError
	}

	grate.Command = command

	unit, err := strconv.Atoi(values[0])
	if err != nil {
		return grate, FormatError
	}

	// limit should > 0
	if limit <= 0 {
		return grate, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
		grate.Period = period
	} else {
		return grate, FormatError
	}

	grate.Limit = limit
	return grate, nil
}
