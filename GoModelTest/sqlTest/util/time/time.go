/*
#Time      :  2021/1/22 2:03 下午
#Author    :  chuangangshen@deepglint.com
#File      :  time.go
#Software  :  GoLand
*/
package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LicTime struct {
	time.Time
}

func (t LicTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(output), nil
}

func (t LicTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LicTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LicTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *LicTime) YearMouthDay() int {
	return 365 * t.Year() + t.YearDay()
}