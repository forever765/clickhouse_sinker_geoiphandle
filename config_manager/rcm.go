package rcm

import (
	"github.com/forever765/clickhouse_sinker_nali/config"
)

// RemoteConfManager can be implemented by many backends: Nacos, Consul, etcd, ZooKeeper...
type RemoteConfManager interface {
	Init(properties map[string]interface{}) error
	GetConfig() (conf *config.Config, err error)
	// PublishConfig publishs the config.
	PublishConfig(conf *config.Config) (err error)
	Register(ip string, port int) (err error)
	Deregister(ip string, port int) (err error)

	// Assignment loop
	Run()
	Stop()
}
