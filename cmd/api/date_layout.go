package main

import (
	"errors"
)

// Time string layouts for parsing.
const (
	Year   = "2006"
	Month  = "2006-01"
	Day    = "2006-01-02"
	Hour   = "2006-01-02 15"
	Minute = "2006-01-02 15:04"
	Second = "2006-01-02 15:04:05"
)

// GetLayout returns the time layout from a string.
func GetLayout(s string) (string, error) {
	var caret, space, points int

	// count relevant characters in the string.
	for _, char := range s {
		switch char {
		case '-':
			caret++
		case ' ':
			space++
		case ':':
			points++
		default:
		}
	}

	// find appropriate layout for parsing.
	switch {
	case caret == 2 && space == 1 && points == 2:
		return Second, nil
	case caret == 2 && space == 1 && points == 1:
		return Minute, nil
	case caret == 2 && space == 1:
		return Hour, nil
	case caret == 2:
		return Day, nil
	case caret == 1:
		return Month, nil
	case caret == 0:
		return Year, nil
	default:
		return "", errors.New("unknown date format " + s)
	}
}
