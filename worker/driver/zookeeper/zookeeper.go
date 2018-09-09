package zookeeper

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/falconbaby/mon-frame/worker"
)

const DriverName = "zookeeper"

type ZookeeperDriver struct{}

func init() {
	worker.Register(DriverName, &ZookeeperDriver{})
}

type ZookeeperConfig struct {
	Step        int64             `json:"step"`
	ConnTimeout int64             `json:"conn_timeout"`
	TranTimeout int64             `json:"tran_timeout"`
	Clusters    map[string]string `json:"cluster"`
}

type zookeeperWorker struct {
	sync.RWMutex
	c *ZookeeperConfig
}

func (d *ZookeeperDriver) Open(config []byte) (worker.Worker, error) {
	var c ZookeeperConfig
	err := json.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}
	return &zookeeperWorker{c: &c}, nil
}

func (w *zookeeperWorker) GetConfig() ([]byte, error) {
	w.RLock()
	defer w.RUnlock()
	return json.Marshal(w.c)
}

func (w *zookeeperWorker) GetStep() int64 {
	w.RLock()
	defer w.RUnlock()
	return w.c.Step
}

func (w *zookeeperWorker) GetTimeout() (conn, tran int64) {
	w.RLock()
	defer w.RUnlock()
	return w.c.ConnTimeout, w.c.TranTimeout
}

func (w *zookeeperWorker) Catch() ([]*worker.Metric, error) {
	fmt.Println("Catch over")
	return nil, nil
}
