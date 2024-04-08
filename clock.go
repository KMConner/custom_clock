package custom_clock

import (
	"context"
	"time"

	"go.uber.org/atomic"
)

type Clock struct {
	speed             *atomic.Float64
	offsetNanoseconds *atomic.Int64
}

func NewClock(speed float64, now Time) *Clock {
	return newClock(speed, now, time.Now())
}

func newClock(speed float64, nowCustom Time, nowReal time.Time) *Clock {
	return &Clock{
		speed:             atomic.NewFloat64(speed),
		offsetNanoseconds: atomic.NewInt64(int64(float64(nowCustom.t.UTC().UnixNano()) - float64(nowReal.UTC().UnixNano())*speed)),
	}
}

func (c *Clock) Speed() float64 {
	return c.speed.Load()
}

func (c *Clock) convertToActualTime(t Time) time.Time {
	return time.Unix(0, int64(float64(t.t.UTC().UnixNano()-c.offsetNanoseconds.Load())/c.speed.Load())).UTC()
}

func (c *Clock) convertFromActualTime(t time.Time) Time {
	return Time{
		t: time.Unix(0, int64(float64(t.UTC().UnixNano())*c.speed.Load())+c.offsetNanoseconds.Load()).UTC(),
	}
}

func (c *Clock) Now() Time {
	return c.convertFromActualTime(time.Now())
}

func (c *Clock) SleepUntil(ctx context.Context, t Time) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(c.calculateRealDuration(time.Now(), t)):
		return nil
	}
}

func (c *Clock) calculateRealDuration(nowReal time.Time, until Time) time.Duration {
	return c.convertToActualTime(until).Sub(nowReal)
}
