package util

import (
	"github.com/robfig/cron/v3"
	"time"
)

func Parse(spec string) (next time.Time, err error) {
	schedule, err := cron.ParseStandard(spec)
	if err != nil {
		return
	}
	next = schedule.Next(time.Now())
	return
}
