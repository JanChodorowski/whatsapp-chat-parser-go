package parser

import (
	"regexp"
	"strconv"
	"strings"
)

// ==================== Time Functions ====================

// convertTime12to24 converts time from 12-hour format to 24-hour format
func convertTime12to24(timeStr string, ampm string) string {
	parts := regexSplitTime.Split(timeStr, -1)
	hours := parts[0]
	minutes := parts[1]
	var seconds string
	if len(parts) > 2 {
		seconds = parts[2]
	}

	if hours == "12" {
		hours = "00"
	}

	hoursInt, _ := strconv.Atoi(hours)
	if ampm == "PM" {
		hoursInt += 12
		hours = strconv.Itoa(hoursInt)
	}

	if seconds != "" {
		return hours + ":" + minutes + ":" + seconds
	}
	return hours + ":" + minutes
}

// normalizeTime normalizes a time string to have the format: hh:mm:ss
func normalizeTime(timeStr string) string {
	parts := regexSplitTime.Split(timeStr, -1)
	hours := parts[0]
	minutes := parts[1]
	seconds := "00"
	if len(parts) > 2 {
		seconds = parts[2]
	}

	// Pad hours with a leading zero if needed
	if len(hours) < 2 {
		hours = "0" + hours
	}

	return hours + ":" + minutes + ":" + seconds
}

// normalizeAMPM normalizes AM/PM indicators to uppercase without other characters
func normalizeAMPM(ampm string) string {
	ampm = regexp.MustCompile(`[^apmAPM]`).ReplaceAllString(ampm, "")
	return strings.ToUpper(ampm)
}
