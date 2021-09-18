package db

import "github.com/housepower/clickhouse_sinker/ipHandle/pkg/dbif"

var dbCache = make(map[dbif.QueryType]dbif.DB)
var queryCache = make(map[string]string)
