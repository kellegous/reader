package internal

import "time"

type Week time.Time

func (w Week) String() string {
	return Day(w).String()
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
