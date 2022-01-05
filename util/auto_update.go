package util

import (
	"github.com/robfig/cron/v3"
	//"../ipHandle/db"
)

// Auto update geoip database file every day
func AddUpdateCronTask() {
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		Logger.Fatal("tick every 1 second")
	})

	c.Start()
}