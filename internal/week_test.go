package internal

import (
	"fmt"
	"testing"
	"time"
)

func mustParseDay(t *testing.T, s string, loc *time.Location) time.Time {
	d, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		t.Error(err)
	}
	return d
}

func TestWeekOf(t *testing.T) {
	tests := []struct {
		Time      time.Time
		Beginning time.Weekday
		Location  *time.Location
		Want      Week
	}{
		{
			Time:      mustParseDay(t, "2025-10-23", time.UTC),
			Beginning: time.Monday,
			Location:  time.UTC,
			Want:      Week(mustParseDay(t, "2025-10-20", time.UTC)),
		},
		{
			Time:      mustParseDay(t, "2025-10-23", time.UTC),
			Beginning: time.Thursday,
			Location:  time.UTC,
			Want:      Week(mustParseDay(t, "2025-10-23", time.UTC)),
		},
		{
			Time:      mustParseDay(t, "2025-10-23", time.UTC),
			Beginning: time.Friday,
			Location:  time.UTC,
			Want:      Week(mustParseDay(t, "2025-10-17", time.UTC)),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s %s", test.Time.Format("2006-01-02"), test.Beginning, test.Location.String()), func(t *testing.T) {
			if got := WeekOf(test.Time, test.Beginning, test.Location); got != test.Want {
				t.Fatalf(
					"WeekOf(%s, %s, %s) = %s, want %s",
					test.Time,
					test.Beginning,
					test.Location.String(),
					got.String(),
					test.Want.String(),
				)
			}
		})
	}
}
