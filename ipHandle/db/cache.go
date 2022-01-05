package db

import "github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/dbif"

var dbCache = make(map[dbif.QueryType]dbif.DB)
var queryCache = make(map[string]string)
