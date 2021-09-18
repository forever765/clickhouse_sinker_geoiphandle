# clickhouse_sinker_geoiphandle

clickhouse_sinker is a sinker program that transfer kafka message into [ClickHouse](https://clickhouse.yandex/).

[Get Started](https://housepower.github.io/clickhouse_sinker/)

Refers to [docs](https://housepower.github.io/clickhouse_sinker/dev/introduction.html#features) to see how it works.  

<br>

## Imagine process
#### Pmacctd --> Kafka --> ClickHouse_Sinker ( ipaddress handle and reduce unknown on class, add "loc_src/loc_dst/isp_src/isp_dst" field )

## Quick Start
configuration new option "geoipHandle" under the "task" field, default value is false

`"geoipHandle": true`