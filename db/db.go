package db

import (
	"sync"
)

type DB struct {
	sync.RWMutex
	data map[string]string
}

func NewDBMap() *DB {
	return &DB{data: make(map[string]string)}
}

func (d *DB) Get(key string) (string, bool) {
	d.RLock()
	defer d.RUnlock()

	val, ok := d.data[key]
	return val, ok
}

func (d *DB) Delete(key string) {
	d.Lock()
	defer d.Unlock()

	delete(d.data, key)
}

func (d *DB) Set(key string, val string) {
	d.Lock()
	defer d.Unlock()

	d.data[key] = val
}
