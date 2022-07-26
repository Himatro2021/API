package helper

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ParseDateAndTimeStringToTime parse date and time in string to time.Time.
// Return error when any error happen
func ParseDateAndTimeStringToTime(date, timeString string) (time.Time, error) {
	dateParts := strings.Split(date, "-")
	timeParts := strings.Split(timeString, ":")
	if len(dateParts) != 3 || len(timeParts) != 2 {
		return time.Time{}, errors.New("date string is invalid")
	}

	year, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return time.Time{}, err
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return time.Time{}, err
	}

	minutes, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return time.Time{}, err
	}

	datetime := time.Date(year, time.Month(month), day, hour, minutes, 0, 0, time.UTC)

	return datetime, nil
}
