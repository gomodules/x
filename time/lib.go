package time

import "time"

// A WeekendAdjustment specifies whether to move before/after/keep the same date.
type WeekendAdjustment int

const (
	NoChange WeekendAdjustment = iota
	Before
	After
)

func AdjustForWeekend(now time.Time, adj WeekendAdjustment) time.Time {
	d := now.Weekday()

	if d == time.Saturday {
		switch adj {
		case NoChange:
			return now
		case Before:
			return now.AddDate(0, 0, -1)
		case After:
			return now.AddDate(0, 0, 2)
		}
	} else if d == time.Sunday {
		switch adj {
		case NoChange:
			return now
		case Before:
			return now.AddDate(0, 0, -2)
		case After:
			return now.AddDate(0, 0, 1)
		}
	}
	return now
}
