package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	timeUtil "github.com/zliang90/kingRest/pkg/util/time"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(timeUtil.LayoutDateTime))
	return []byte(output), nil
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
