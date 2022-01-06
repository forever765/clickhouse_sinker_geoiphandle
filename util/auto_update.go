package util

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"

	"github.com/forever765/clickhouse_sinker_nali/ipHandle/db"
)

// Auto update geoip database file every day
func AddUpdateCronTask(JobInterval string) {
	c := cron.New()
	c.AddFunc(JobInterval, DoUpdate)
	c.Start()
	Logger.Info("Auto update cron job running, time interval: ", zap.String("",JobInterval))
	time.After(time.Hour * 168)
}

func DoUpdate() {
	db.Update()
	Logger.Info("tick every 5 second")
}