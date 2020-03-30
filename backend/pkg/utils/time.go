package utils

import (
	"os"
	"time"
)

func TestNowUTC() time.Time {
	t, _ := time.Parse(time.UnixDate, "Sat Mar  7 11:06:39 PST 2015")

	return t.UTC()
}

func NowUTC() time.Time {
	if os.Getenv("NODE_ENV") == "test" {
		return TestNowUTC()
	}

	return time.Now().UTC()
}
