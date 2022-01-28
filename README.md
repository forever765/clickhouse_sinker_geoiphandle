# Clickhouse_Sinker_Nali

Clickhouse_Sinker is a sinker program that transfer kafka message into [ClickHouse](https://clickhouse.yandex/).
Clickhouse_Sinker_Nali 

[clickhouse_sinker docs](https://housepower.github.io/clickhouse_sinker_nali/dev/introduction.html#features)  

<br>

- #### Clickhouse_Sinker_Nali based on [Clickhouse_Sinker v2.2.0](https://github.com/forever765/clickhouse_sinker)
- #### GeoIP information provide from [Nali v0.3.5](https://github.com/zu1k/nali)
- #### Import robfig/cron/v3 package to auto update geoip database file every day

## Processing flow
##### Pmacctd --> Kafka --> ClickHouse_Sinker_nali --> Clickhouse
1. Get messages from Kafka
2. Get ip_src and ip_dst geo info from Nali module
3. Reduce unknown on class field (serviceMap replace unknown)
4. Add "loc_src/loc_dst/isp_src/isp_dst" field to message
5. Write messages to Clickhouse

## Build && Run
`go get -u github.com/forever765/clickhouse_sinker_nali/...`
`make build`

## Quick Start
configuration new option "geoipHandle" & "autoUpdateGeoIPDB" under the "task" field.  
1. geoipHandle: use or not handle geoip info.
2. autoUpdateGeoIPDB: set auto update NALI geoip db file cronjob interval time.

`"geoipHandle": true`  
`"autoUpdateGeoIPDB": "0 3 * * *"`

## Note
1. Sinker listen port and log path setting: `cmdOps on cmd/clickhouse_sinker_nali/main.go`
2. GeoIP Database file download path: ` variable "HomePath" on ipHandle/constant/path.go`