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

func (rs Rates) Append(sr singleRate) {
	rs = append(rs, sr)
}

var timeDict = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": time.Hour * 24,
}

var CommandError = errors.New("The command of first number should > 0.")
var FormatError = errors.New("Please check the format with your input.")
var LimitError = errors.New("Limit should > 0.")

// NewGlobalRate("10-m", 200)
// Each 10 minutes single ip address can request 200 times.
func newGlobalRate(command string, limit int) (GlobalRate, error) {
	var gRate GlobalRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return gRate, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	log.Printf("unit = %d", unit)
	if err != nil {
		return gRate, FormatError
	}
	if unit <= 0 {
		return gRate, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return gRate, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return gRate, FormatError
	}

	gRate.Command = command
	gRate.Period = period
	gRate.Limit = limit
	return gRate, nil
}

func newSingleRate(path, command string, limit int) (singleRate, error) {
	var sRate singleRate
	var period time.Duration

	values := strings.Split(command, "-")
	if len(values) != 2 {
		log.Println("Some error with your input command!, the len of your command is ", len(values))
		return sRate, FormatError
	}

	unit, err := strconv.Atoi(values[0])
	log.Printf("unit = %d", unit)
	if err != nil {
		return sRate, FormatError
	}
	if unit <= 0 {
		return sRate, CommandError
	}

	// limit should > 0
	if limit <= 0 {
		return sRate, LimitError
	}

	if t, ok := timeDict[strings.ToUpper(values[1])]; ok {
		period = time.Duration(unit) * t
	} else {
		return sRate, FormatError
	}

	sRate.Path = path
	sRate.Command = command
	sRate.Period = period
	sRate.Limit = limit
	return sRate, nil
}
