package main

import "time"

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

type JsonResult struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func newJsonResult(code int, message string, data interface{}) JsonResult {
	return JsonResult{
		Code:code,
		Message: message,
		Data: data,
	}
}

func newSuccessJsonResult(message string, data interface{}) JsonResult {
	return JsonResult{
		Code: 0,
		Message: message,
		Data: data,
	}
}

func newErrorJsonResult(message string) JsonResult {
	return JsonResult{
		Code: 1,
		Message: message,
	}
}

const SessionName = "rad_access_token"

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}

func NowTime() Time {
	return Time(time.Now())
}