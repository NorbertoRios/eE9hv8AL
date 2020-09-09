package utils

import (
	"regexp"
	"time"
)

//MinTimeStamp generates min timestamp value
func MinTimeStamp() time.Time {
	min, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	return min
}

var dateRegex, _ = regexp.Compile("^\\d{4}\\/\\d{2}\\/\\d{2}$")

//IsValidStringDate checks date part of string for mask 2000/01/01
func IsValidStringDate(date string) bool {
	return dateRegex.Match([]byte(date))
}

var timeRegex, _ = regexp.Compile("^\\d{2}\\:\\d{2}\\:\\d{2}$")

//IsValidStringTime checks time part of string for mask 00:00:00
func IsValidStringTime(time string) bool {

	return timeRegex.Match([]byte(time))
}
