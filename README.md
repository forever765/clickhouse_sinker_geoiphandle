# clickhouse_sinker_nali

clickhouse_sinker is a sinker program that transfer kafka message into [ClickHouse](https://clickhouse.yandex/).

[clickhouse_sinker docs](https://housepower.github.io/clickhouse_sinker/dev/introduction.html#features)  

<br>

- clickhouse_sinker_nali base clickhouse_sinker v1.91

- #### GeoIP information provide from [nali](https://github.com/zu1k/nali)

## Processing flow
##### Pmacctd --> Kafka --> ClickHouse_Sinker_nali
1. Get ip_src and ip_dst geo info
2. Reduce unknown on class field
3. Add "loc_src/loc_dst/isp_src/isp_dst" field to kafka
4. If enum8 or enum16 type field exists in ClickHouse, ck_sinker will treat it as a string type

## Build && Run
`make build`

## Quick Start
configuration new option "geoipHandle" under the "task" field, default value is false

`"geoipHandle": true`
