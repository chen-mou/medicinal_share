package entity

import (
	"database/sql/driver"
	"time"
)

type Time struct {
	tim *time.Time
}

type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func CreatePage(num int, size int) *Page {
	return &Page{
		Limit:  size,
		Offset: size * (num - 1),
	}
}

func (t *Time) Scan(value interface{}) error {
	time := value.(time.Time)
	t.tim = &time
	return nil
}

func (t Time) Value() (driver.Value, error) {
	if t.tim == nil {
		return nil, nil
	}
	return *t.tim, nil
}

func (t Time) Time() time.Time {
	return *t.tim
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
