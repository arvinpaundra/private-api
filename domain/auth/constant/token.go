package constant

import "time"

const (
	TokenValidSevenDays       = 7 * 24 * time.Hour
	TokenValidThreeHours      = 3 * time.Hour
	TokenValidAfterThreeHours = 3 * time.Hour
	TokenValidImmediately     = 0 * time.Minute
)
