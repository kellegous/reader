package internal

import (
	"fmt"
	"time"
)

type Week time.Time

func (w Week) BeginsAt() time.Time {
	return Day(w).AsTime()
}

func (w Week) EndsAt() time.Time {
	return w.BeginsAt().AddDate(0, 0, 7).Add(-1 * time.Nanosecond)
}

func (w Week) String() string {
	return fmt.Sprintf("%s - %s", Day(w.BeginsAt()), Day(w.EndsAt()))
}

func WeekOf(
	t time.Time,
	beginning time.Weekday,
	loc *time.Location,
) Week {
	d := DayOf(t, loc)
	offset := int(beginning) - int(t.Weekday())
	if offset > 0 {
		// earlier week
		d = d.Add(-7 + offset)
	} else if offset < 0 {
		// later week
		d = d.Add(offset)
	}
	return Week(d)
}
