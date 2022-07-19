package helper

import "time"

// IsStartAndFinishTimeValid validate whether the start and finish time
// valid, and are they placed in order (start must happen before finish in time)
func IsStartAndFinishTimeValid(start, finish time.Time) bool {
	if start.IsZero() || finish.IsZero() {
		return false
	}

	if start.After(finish) || finish.Before(start) {
		return false
	}

	now := time.Now()

	return !finish.Before(now)
}
