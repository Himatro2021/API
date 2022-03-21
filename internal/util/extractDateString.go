package util

import (
	"errors"
	"strconv"
	"strings"
)

func ExtractDateString(dateString string) (int, int, int, error) {
	splited := strings.Split(dateString, "-")

	if len(splited) != 3 {
		return 0, 0, 0, errors.New("invalid date string")
	}

	year, err := strconv.Atoi(splited[0])

	if err != nil {
		return 0, 0, 0, errors.New("invalid date string")
	}

	month, err := strconv.Atoi(splited[1])

	if err != nil {
		return 0, 0, 0, errors.New("invalid date string")
	}

	if month < 1 || month > 12 {
		return 0, 0, 0, errors.New("invalid month value in date string")
	}

	day, err := strconv.Atoi(splited[2])

	if err != nil {
		return 0, 0, 0, errors.New("invalid date string")
	}

	if day < 1 || day > 31 {
		return 0, 0, 0, errors.New("invalid day value in date string")
	}

	return year, month, day, nil
}
