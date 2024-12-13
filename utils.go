package main

import (
	"regexp"
	"strconv"
	"time"
)

func reExtract(pattern string, str string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractCSRFToken(html string) string {
	return reExtract(
		`<input[^>]+id=\"csrftoken\"[^>]+value=\"([^\"]+)\"`,
		html)
}

func extractRTK(js string) string {
	return reExtract(`rtk:'([a-f0-9-]+)'`, js)
}

func getTimestamp() int64 {
	return time.Now().UnixMilli()
}

func getTimestampStr() string {
	return strconv.FormatInt(getTimestamp(), 10)
}
