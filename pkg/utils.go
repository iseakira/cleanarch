package pkg

import (
	"os"
	"time"
)

// Str2times Converts the string t passed in YYYY-MM-DD format to time.Time type and returns it.
func Str2time(t string) time.Time {
	parsedTime, _ := time.Parse("2006-01-02", t)
	return parsedTime
}

func GetEnvDefault(key, defVal string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defVal
	}
	return val
}