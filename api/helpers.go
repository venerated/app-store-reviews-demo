package main

import (
	"time"
)

func convertTimestamp(timestamp string, inputFormat string, outputFormat string) string {
	parsedTime, err := time.Parse(inputFormat, timestamp)
	if err != nil {
		return ""
	}
	return parsedTime.Format(outputFormat)
}
