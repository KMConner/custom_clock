package custom_clock

import (
	"time"
)

type Time struct {
	t time.Time // actual time
}

func NewCustomTime(t time.Time) Time {
	return Time{
		t: t.UTC(),
	}
}

func (t Time) Format(format string) string {
	loc := time.FixedZone("CUSTOM", 0)
	tt := t.t
	return tt.In(loc).Format(format)
}

func (t Time) Sub(t1 Time) time.Duration {
	return t.t.Sub(t1.t)
}
