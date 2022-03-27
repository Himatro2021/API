package util

import (
	"errors"
	"strconv"
	"strings"
)

func ExtractTimeString(timeString string) (int, int, int, error) {
	splited := strings.Split(timeString, ":")

	if len(splited) == 2 {
		splited = append(splited, "00") // if user send only hour and minute
	}

	if len(splited) != 3 {
		LogErr("WARN", "invalid date string", "")
		return 0, 0, 0, errors.New("invalid date string")
	}

	hour, err := strconv.Atoi(splited[0])

	if err != nil {
		LogErr("WARN", "invalid date string", "")
		return 0, 0, 0, errors.New("invalid date string")
	}

	minute, err := strconv.Atoi(splited[1])

	if err != nil {
		LogErr("WARN", "invalid date string", "")
		return 0, 0, 0, errors.New("invalid date string")
	}

	sec, err := strconv.Atoi(splited[2])

	if err != nil {
		LogErr("WARN", "invalid date string", "")
		return 0, 0, 0, errors.New("invalid date string")
	}

	return hour, minute, sec, nil
}
