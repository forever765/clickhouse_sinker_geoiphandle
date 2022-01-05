package util

import (
	"github.com/robfig/cron/v3"
	//"../ipHandle/db"
)

// Auto update geoip database file every day
func AddUpdateCronTask(JobInterval string) {
	c := cron.New()
	c.AddFunc(JobInterval, haha)
	c.Start()
	select{}
}

func haha() {
	Logger.Fatal("tick every 5 second")
}