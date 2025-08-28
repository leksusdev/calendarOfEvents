package datetime

import (
	"strings"
	"time"
)

const LayoutFormat = "2006-01-02 15:04"

func ParseLocal(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	return time.ParseInLocation(LayoutFormat, s, time.Local)
}

func NormalizeUTCSeconds(t time.Time) time.Time {
	return t.UTC().Truncate(time.Second)
}

func FormatLocal(t time.Time) string {
	return t.In(time.Local).Format(LayoutFormat)
}
