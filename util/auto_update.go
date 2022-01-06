package util

import (
	"fmt"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/constant"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/cdn"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/qqwry"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/zxipv6wry"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

var (
	QQWryPath = filepath.Join(constant.HomePath, "qqwry.dat")
	ZXIPv6WryPath = filepath.Join(constant.HomePath, "zxipv6wry.db")
	CDNPath = filepath.Join(constant.HomePath, "cdn.json")
)

// Auto update geoip database file every day
func AddUpdateCronTask(JobInterval string) {
	c := cron.New()
	c.AddFunc(JobInterval, DoUpdate)
	c.Start()
	Logger.Info("Add cron job: auto update geoip db file succeed, time interval: ", zap.String("",JobInterval))
	time.After(time.Hour * 168)
}

func DoUpdate() {
	fmt.Println("haha")
	startTime := time.Now().UnixNano()
	qqwry.Download(QQWryPath)
	zxipv6wry.Download(ZXIPv6WryPath)
	cdn.Download(CDNPath)
	endTime := time.Now().UnixNano()
	Logger.Info("Update geoip db done, Elapsed time: ", zap.Float64("", float64(endTime-startTime)/1000000))
}