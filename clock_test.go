package custom_clock

import (
	"testing"
	"time"
)

func TestClockConvert(t *testing.T) {
	t.Parallel()
	time1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2021, 1, 1, 1, 0, 1, 0, time.UTC)
	time3 := time.Date(2021, 1, 1, 2, 0, 2, 0, time.UTC).In(time.FixedZone("CUSTOM", 9*3600))
	time4 := time.Date(2021, 1, 1, 3, 0, 3, 0, time.UTC)

	testCases := []struct {
		name     string
		speed    float64
		nowFake  Time
		nowReal  time.Time
		realTime time.Time
		fakeTime Time
	}{
		{
			name:     "same as actual",
			speed:    1,
			nowFake:  NewCustomTime(time1),
			nowReal:  time1,
			realTime: time2,
			fakeTime: NewCustomTime(time2),
		},
		{
			name:     "1 hour forward then real but same speed",
			speed:    1,
			nowFake:  NewCustomTime(time2),
			nowReal:  time1,
			realTime: time2,
			fakeTime: NewCustomTime(time3),
		},
		{
			name:     "1 hour behind then real but same speed",
			speed:    1,
			nowFake:  NewCustomTime(time1),
			nowReal:  time2,
			realTime: time3,
			fakeTime: NewCustomTime(time2),
		},
		{
			name:     "no offset but 2x speed",
			speed:    2,
			nowFake:  NewCustomTime(time1),
			nowReal:  time1,
			realTime: time2,
			fakeTime: NewCustomTime(time3),
		},
		{
			name:     "no offset but 0.5x speed",
			speed:    0.5,
			nowFake:  NewCustomTime(time1),
			nowReal:  time1,
			realTime: time3,
			fakeTime: NewCustomTime(time2),
		},
		{
			name:     "1 hour forward then real and 2x speed",
			speed:    2,
			nowFake:  NewCustomTime(time2),
			nowReal:  time1,
			realTime: time2,
			fakeTime: NewCustomTime(time4),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			clock := newClock(tc.speed, tc.nowFake, tc.nowReal)
			actualTime := clock.convertToActualTime(tc.fakeTime)
			if !actualTime.Equal(tc.realTime) {
				t.Errorf("actualTime: %v, realTime: %v", actualTime, tc.realTime)
			}
			fakeTime := clock.convertFromActualTime(tc.realTime)
			if !fakeTime.t.Equal(tc.fakeTime.t) {
				t.Errorf("fakeTime: %v, fakeTime: %v", fakeTime.t, tc.fakeTime.t)
			}
		})
	}
}

func TestCalculateRealDuration(t *testing.T) {
	t.Parallel()
	time1 := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2021, 1, 1, 1, 0, 1, 0, time.UTC)
	time3 := time.Date(2021, 1, 1, 2, 0, 2, 0, time.UTC).In(time.FixedZone("CUSTOM", 9*3600))
	time4 := time.Date(2021, 1, 1, 3, 0, 3, 0, time.UTC)

	testCases := []struct {
		name       string
		speed      float64
		initCustom Time
		initReal   time.Time
		until      Time
		nowReal    time.Time
		want       time.Duration
	}{
		{
			name:       "same as actual, no duration",
			speed:      1,
			initCustom: NewCustomTime(time1),
			initReal:   time1,
			nowReal:    time2,
			until:      NewCustomTime(time2),
			want:       0,
		},
		{
			name:       "same as actual",
			speed:      1,
			initCustom: NewCustomTime(time1),
			initReal:   time1,
			nowReal:    time2,
			until:      NewCustomTime(time3),
			want:       time.Hour + time.Second,
		},
		{
			name:       "1 hour forward then real but same speed",
			speed:      1,
			initCustom: NewCustomTime(time2),
			initReal:   time1,
			nowReal:    time3,
			until:      NewCustomTime(time3),
			want:       -time.Hour - time.Second,
		},
		{
			name:       "1 hour forward then real but same speed 2",
			speed:      1,
			initCustom: NewCustomTime(time2),
			initReal:   time1,
			nowReal:    time3,
			until:      NewCustomTime(time2),
			want:       -2*time.Hour - 2*time.Second,
		},
		{
			name:       "no offset but 2x speed",
			speed:      2,
			initCustom: NewCustomTime(time1),
			initReal:   time1,
			nowReal:    time2,
			until:      NewCustomTime(time3),
			want:       0,
		},
		{
			name:       "no offset but 2x speed 2",
			speed:      2,
			initCustom: NewCustomTime(time1),
			initReal:   time1,
			nowReal:    time2,
			until:      NewCustomTime(time4),
			want:       time.Hour/2 + time.Second/2,
		},
		{
			name:       "1 hour forward then real and 2x speed",
			speed:      2,
			initCustom: NewCustomTime(time2),
			initReal:   time1,
			nowReal:    time3,
			until:      NewCustomTime(time4),
			want:       -time.Hour - time.Second,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			clock := newClock(tc.speed, tc.initCustom, tc.initReal)
			got := clock.calculateRealDuration(tc.nowReal, tc.until)
			if tc.want != got {
				t.Errorf("got: %v, want: %v", got, tc.want)
			}
		})
	}
}
