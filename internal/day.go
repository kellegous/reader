package internal

import "time"

type Day time.Time

func DayOf(t time.Time, loc *time.Location) Day {
	t = t.In(loc)
	return Day(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc))
}

func (d Day) Add(days int) Day {
	return Day(time.Time(d).AddDate(0, 0, days))
}

func (d Day) String() string {
	return time.Time(d).Format("2006-01-02")
}
