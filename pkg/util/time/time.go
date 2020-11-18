package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Layout = string

const (
	// year-month-day
	LayoutDate Layout = "2006-01-02"

	// year-month-day hour:minute:second
	LayoutDateTime Layout = "2006-01-02 15:04:05"

	// hour-minute
	LayoutTime = "15:04"
)

// NowString format time now
func NowString(layout Layout) string {
	return time.Now().Format(layout)
}

// HourMinute hour and minute
type HourMinute struct {
	Hour   int
	Minute int
}

func (hm HourMinute) String() string {
	h := fmt.Sprintf("%d", hm.Hour)
	if hm.Hour < 10 {
		h = fmt.Sprintf("0%s", h)
	}
	m := fmt.Sprintf("%d", hm.Minute)
	if hm.Minute < 10 {
		m = fmt.Sprintf("0%s", m)
	}
	return fmt.Sprintf("%s:%s", h, m)
}

// ParseTimeString parse time string, convert "09:00" to HourMinute type
func ParseTimeString(hm string) (hourMinute *HourMinute, err error) {
	defer func() {
		if e := recover(); e != nil {
			hourMinute = nil
			err = fmt.Errorf("can not parse hourMinute string, err: %v", e)
		}
	}()

	if hm != "" && strings.Contains(hm, ":") {
		hmList := strings.Split(hm, ":")

		hour, herr := strconv.Atoi(hmList[0])
		if herr != nil {
			return nil, herr
		}
		min, minErr := strconv.Atoi(hmList[1])
		if minErr != nil {
			return nil, minErr
		}
		return &HourMinute{hour, min}, nil
	}
	return nil, fmt.Errorf("illegal hourMinute format, eg: '08:30'")
}

// ToTimeString convert time hour or minute to string,
func ParseTimeInt(i int) string {
	if i < 10 {
		return fmt.Sprintf("0%d", i)
	}
	return fmt.Sprintf("%d", i)
}

// ParseHourMinute convert time string to HourMinute type
func ParseHourMinuteString(hm string) (hour, minute string, err error) {
	if !(strings.Contains(hm, ":") && len(hm) == 5) {
		err = fmt.Errorf("illegal patrol_start_time: %s, eg: '00:30'", hm)
		return
	}

	var h, m int

	t := strings.Split(hm, ":")
	if h, err = strconv.Atoi(t[0]); err != nil {
		return
	}
	if m, err = strconv.Atoi(t[1]); err != nil {
		return
	}
	if h < 0 || h > 23 {
		err = fmt.Errorf("illegal time hour: '%s', eg: '00 ~ 23'", t[0])
		return
	}
	if m < 0 || m > 59 {
		err = fmt.Errorf("illegal time minute: '%s', eg: '00 ~ 59'", t[1])
	}

	hour = ParseTimeInt(h)
	minute = ParseTimeInt(m)
	return
}

// TimeStampInScope 时间戳是否在指定时间范围内
func TimeStampInScope(timestamp int64, d time.Duration) bool {
	return time.Now().Unix()-timestamp <= int64(d.Seconds())
}
