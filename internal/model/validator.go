package model

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	validate "github.com/go-playground/validator/v10"
)

// Validator singleton for validating struct
var Validator = validate.New()

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		_ = Validator.RegisterValidation("time", isTimeFormatValid)
		_ = Validator.RegisterValidation("date", isDateFormatValid)
	})
}

func isTimeFormatValid(fl validate.FieldLevel) bool {
	timeString := fl.Field().String()
	if timeString == "" {
		return false
	}

	parts := strings.Split(timeString, ":")
	if len(parts) != 2 {
		return false
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}

	if hour < 0 || hour > 23 {
		return false
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	if minute < 0 || minute > 59 {
		return false
	}

	return true
}

func isDateFormatValid(fl validate.FieldLevel) bool {
	dateString := fl.Field().String()

	if dateString == "" {
		fmt.Println("0")
		return false
	}

	parts := strings.Split(dateString, "-")
	if len(parts) != 3 {
		return false
	}

	if err := datePartsValidator(parts); err != nil {
		return false
	}

	_, err := time.Parse("2006-01-02", dateString)
	return err == nil
}

func datePartsValidator(parts []string) (err error) {
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return
	}

	if year < 1 || month < 1 || day < 1 {
		return
	}

	if month > 12 || day > 31 {
		return
	}

	return
}
