package main

import (
	"time"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
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

func newJsonResult(code int, message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func newSuccessJsonResult(message string, data interface{}) JsonResult {
	return JsonResult{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

func defaultSuccessJsonResult(data interface{}) JsonResult {
	return JsonResult{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func newErrorJsonResult(message string) JsonResult {
	return JsonResult{
		Code:    1,
		Message: message,
	}
}

const SessionName = "rad_access_token"

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	if s == "" || s == "null" {
		// return zero value
		*t = Time(time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))
		return
	}
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, s, time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	// when time is zero value return empty string
	if time.Time(t).IsZero() {
		return []byte(`""`), nil
	}

	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}

func (t *Time) convert(datetime time.Time) Time {
	return Time(datetime)
}

func NowTime() Time {
	return Time(time.Now())
}

func getStdTimeFromString(value string) (time.Time, error) {
	return time.ParseInLocation(TimeFormat, value, time.Local)
}

func getTodayLastTime() time.Time {
	today := time.Now()
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		23, 59, 59, 0, today.Location()))
}

func getNextDayLastTime() time.Time {
	today := time.Now()
	return time.Time(time.Date(today.Year(), today.Month(), today.Day()+1,
		23, 59, 59, 0, today.Location()))
}

func getMonthLastTime() time.Time {
	today := time.Now()
	today = today.AddDate(0, 1, -today.Day())
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		23, 59, 59, 0, today.Location()))
}

func getDayLastTimeAfterAYear() time.Time {
	today := time.Now()
	today = today.AddDate(1, 0, 0)
	return time.Time(time.Date(today.Year(), today.Month(), today.Day(),
		23, 59, 59, 0, today.Location()))
}
