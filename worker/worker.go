package worker

import (
	"fmt"
	"sort"
	"sync"
)

var (
	locker  sync.RWMutex
	drivers = make(map[string]Driver)
)

type Metric struct{}

type Driver interface {
	Open(config []byte) (Worker, error)
}

type Worker interface {
	GetConfig() ([]byte, error)
	GetTimeout() (conn, tran int64)
	GetStep() int64
	Catch() ([]*Metric, error)
}

func Register(name string, driver Driver) {
	locker.Lock()
	defer locker.Unlock()
	if driver == nil {
		panic("worker: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("worker: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Drivers() []string {
	locker.RLock()
	defer locker.RUnlock()
	var list = make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func Open(name string, config []byte) (Worker, error) {
	locker.RLock()
	driver, ok := drivers[name]
	locker.RUnlock()
	if !ok {
		return nil, fmt.Errorf("worker: unknown driver %q (forgotten import?)", name)
	}
	return driver.Open(config)
}
