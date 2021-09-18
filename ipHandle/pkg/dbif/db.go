package dbif

import (
	"fmt"

	"github.com/housepower/clickhouse_sinker/ipHandle/pkg/cdn"
	"github.com/housepower/clickhouse_sinker/ipHandle/pkg/geoip"
	"github.com/housepower/clickhouse_sinker/ipHandle/pkg/ipip"
	"github.com/housepower/clickhouse_sinker/ipHandle/pkg/qqwry"
	"github.com/housepower/clickhouse_sinker/ipHandle/pkg/zxipv6wry"
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
