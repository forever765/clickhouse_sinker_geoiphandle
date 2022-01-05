# Run

## Requirements

Note: Ensure `clickhouse-server` and `kafka` work before running clickhouse_sinker.

## Configs

> There are two ways to get config: a local single config, or Nacos.

- For local file:

  `clickhouse_sinker --local-cfg-file docker/test_auto_schema.json`

- For Nacos:

  `clickhouse_sinker --nacos-addr 127.0.0.1:8848 --nacos-username nacos --nacos-password nacos --nacos-dataid test_auto_schema`

> Read more detail descriptions of config in [here](../configuration/config.html)

## Example

Let's follow up a piece of the systest script.

* Prepare

  - let's checkout `clickhouse_sinker`

  ```bash
  $ git clone https://github.com/forever765/clickhouse_sinker_nali.git
  $ cd clickhouse_sinker_nali
  ```

  - let's start standalone clickhouse-server and kafka in container:

  ```bash
  $ docker-compose up -d
  ```

* Create a simple table in Clickhouse

  > It's not the duty for clickhouse_sinker to auto create table, so we should do that manually.

  ```sql
  CREATE TABLE IF NOT EXISTS test_auto_schema
  (
      `day` Date DEFAULT toDate(time),
      `time` DateTime,
      `name` String,
      `value` Float64
  )
  ENGINE = MergeTree
  PARTITION BY day
  ORDER BY (time, name);
  ```

* Create a topic in kafka

  > I use [kaf](https://github.com/birdayz/kaf) tool to create topics.

  ```bash
  $ kaf topic create topic1 -p 1 -r 1
  ✅ Created topic!
        Topic Name:            topic1
        Partitions:            1
        Replication Factor:    1
        Cleanup Policy:        delete
  ```


* Run clickhouse_sinker

  ```bash
  $ ./clickhouse_sinker_nali --local-cfg-file docker/test_auto_schema.json
  ```


* Send messages to the topic

  ```bash
  echo '{"time" : "2020-12-18T03:38:39.000Z", "name" : "name1", "value" : 1}' | kaf -b '127.0.0.1:9092' produce topic1
  echo '{"time" : "2020-12-18T03:38:39.000Z", "name" : "name2", "value" : 2}' | kaf -b '127.0.0.1:9092' produce topic1
  echo '{"time" : "2020-12-18T03:38:39.000Z", "name" : "name3", "value" : 3}' | kaf -b '127.0.0.1:9092' produce topic1
  ```

* Check the data in clickhouse

  ```sql
  SELECT count() FROM test_auto_schema;

  3 rows in set. Elapsed: 0.016 sec.

  ```
