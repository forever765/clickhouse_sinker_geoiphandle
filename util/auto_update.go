package util

import (
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/constant"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"os/exec"
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
	Logger.Info("Add cron job: auto update geoip db file succeed, ", zap.String("time interval: ",JobInterval))
	time.After(time.Hour * 168)
}

func DoUpdate() {
	Logger.Info("Start update Geoip database...")
	startTime := time.Now().UnixNano()
	QqwryDownload(QQWryPath)
	Zxipv6wry_Download(ZXIPv6WryPath)
	CdnDownload(CDNPath)
	endTime := time.Now().UnixNano()
	Logger.Info("Update Geoip database file done, ", zap.Float64("Elapsed time (second):", float64(endTime-startTime)/1000000000))
	Logger.Info("Restarting_Clickhouse Sinker_Nali")
	restartMyself()

}

func restartMyself(){
	cmd := exec.Command("/usr/bin/systemctl", "restart", "ch_sinker")
	cmd.Run()
}