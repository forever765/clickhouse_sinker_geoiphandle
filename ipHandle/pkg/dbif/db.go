package dbif

import (
	"fmt"

	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/cdn"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/geoip"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/ipip"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/qqwry"
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/zxipv6wry"
)

type QueryType uint

const (
	TypeIPv4 = iota
	TypeIPv6
	TypeDomain
)

type DB interface {
	Find(query string, params ...string) (result fmt.Stringer, err error)
}

var (
	_ DB = qqwry.QQwry{}
	_ DB = zxipv6wry.ZXwry{}
	_ DB = ipip.IPIPFree{}
	_ DB = geoip.GeoIP{}
	_ DB = cdn.CDN{}
)
