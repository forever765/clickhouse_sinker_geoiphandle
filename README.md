# clickhouse_sinker_geoiphandle

clickhouse_sinker is a sinker program that transfer kafka message into [ClickHouse](https://clickhouse.yandex/).

[Get Started](https://housepower.github.io/clickhouse_sinker/)

Refers to [docs](https://housepower.github.io/clickhouse_sinker/dev/introduction.html#features) to see how it works.  

<br>

## Imagine process
##### Pmacctd --> Kafka --> ClickHouse_Sinker(ipaddress handle, add "ip_src_country" and "ip_dst_country" field)

## Still under development