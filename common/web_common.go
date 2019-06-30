package common

import (
	"time"
)

const (
	DateFormat  = "2006-01-02"
	TimeFormat  = "2006-01-02 15:04:05"
	Zore        = 0
	One         = 1
	Ten         = 10
	FiftyNine   = 59
	TwentyThree = 23
)

const (
	MonthlyProduct = 1
	TimeProduct    = 2
	FlowProduct    = 3
)

const (
	DefaultFlowClearCycle     = 1
	DayFlowClearCycle         = 2
	MonthFlowClearCycle       = 3
	FixedPeriodFlowClearCycle = 4
)

type JsonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewJsonResult(code int, message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewSuccessJsonResult(message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    Zore,
		Message: message,
		Data:    data,
	}
}

func DefaultSuccessJsonResult(data interface{}) JsonResult {
	return JsonResult{
		Code:    Zore,
		Message: "success",
		Data:    data,
	}
}

func NewErrorJsonResult(message string) JsonResult {
	return JsonResult{
		Code:    One,
		Message: message,
	}
}

const SessionName = "Authorization"

func GetStdTimeFromString(value string) (time.Time, error) {
	return time.ParseInLocation(TimeFormat, value, time.Local)
}

func GetTodayLastTime() time.Time {
	today := time.Now()
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		TwentyThree, FiftyNine, FiftyNine, Zore, today.Location()))
}

func GetNextDayLastTime() time.Time {
	today := time.Now()
	return time.Time(time.Date(today.Year(), today.Month(), today.Day()+1,
		TwentyThree, FiftyNine, FiftyNine, Zore, today.Location()))
}

func GetMonthLastTime() time.Time {
	today := time.Now()
	today = today.AddDate(Zore, One, -today.Day())
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		TwentyThree, FiftyNine, FiftyNine, Zore, today.Location()))
}

func GetDayLastTimeAfterAYear() time.Time {
	today := time.Now()
	today = today.AddDate(One, Zore, Zore)
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		TwentyThree, FiftyNine, FiftyNine, Zore, today.Location()))
}
