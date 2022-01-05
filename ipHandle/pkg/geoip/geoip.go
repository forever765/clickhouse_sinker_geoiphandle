package geoip

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/forever765/clickhouse_sinker_nali/util"
	"github.com/oschwald/geoip2-golang"
)

// GeoIP2
type GeoIP struct {
	db *geoip2.Reader
}

// new geoip from database file
func NewGeoIP(filePath string) (geoip GeoIP) {
	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 Geoip2 City库，并保存在", filePath)
		os.Exit(1)
	} else {
		db, err := geoip2.Open(filePath)
		if err != nil {
			log.Fatal(err)
		} else {
			util.Logger.Info("maxmind db 已加载！")
		}
		geoip = GeoIP{db: db}
	}
	return
}

func (g GeoIP) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be valid IP")
	}
	record, err := g.db.City(ip)
	if err != nil {
		return
	}

	lang := "zh-CN"
	if len(params) > 0 {
		if _, ok := record.Country.Names[params[0]]; ok {
			lang = params[0]
		}
	}

	result = Result{
		Country: record.Country.Names[lang],
		City:    record.City.Names[lang],
	}
	return
}

type Result struct {
	Country string
	City    string
}

func (r Result) String() string {
	if r.City == "" {
		return r.Country
	} else {
		return fmt.Sprintf("%s %s", r.Country, r.City)
	}
}
