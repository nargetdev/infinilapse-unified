package stopgap

import (
	"fmt"
	"time"
)

const (
	LIGHTS_ON  = "08:00"
	LIGHTS_OFF = "23:59"
	LOC        = "MST"
)

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func STOPGAP_getFadeByTime() float64 {
	newLayout := "15:04"
	loc, _ := time.LoadLocation(LOC)
	check := time.Now().In(loc)
	start, _ := time.ParseInLocation(newLayout, LIGHTS_ON, loc)
	end, _ := time.ParseInLocation(newLayout, LIGHTS_OFF, loc)
	start = start.AddDate(check.Year(), int(check.Month())-1, check.Day()-1)
	end = end.AddDate(check.Year(), int(check.Month())-1, check.Day()-1)
	isDaytime := inTimeSpan(start, end, check)
	fmt.Println(start.String()+" --- "+check.String()+" --- ", end.String(), isDaytime)

	if isDaytime {
		return 1.0
	} else {
		return 0.0
	}
}
