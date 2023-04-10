package entity

import (
	"database/sql/driver"
	"time"
)

type Time struct {
	tim *time.Time
}

func (t *Time) Scan(value interface{}) error {
	*t.tim = value.(time.Time)
	return nil
}

func (t Time) Value() (driver.Value, error) {
	if t.tim == nil {
		return nil, nil
	}
	return *t.tim, nil
}

func Now() Time {
	now := time.Now()
	return Time{
		tim: &now,
	}
}

func CreateTime(t time.Time) Time {
	return Time{
		&t,
	}
}

func (t *Time) MarshalJSON() ([]byte, error) {
	s := t.tim.Format("2006-01-02 15-04-05")
	return []byte("\"" + s + "\""), nil
}

func (t *Time) UnmarshalJSON(byt []byte) error {
	tim := string(byt)
	t1, _ := time.Parse("\"2006-01-02 15-04-05\"", tim)
	t.tim = &t1
	return nil
}
