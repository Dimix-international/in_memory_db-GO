package db

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	data map[string]string
}

type DB struct {
	shardMap []*Shard
}

func NewShardMap(count int) *DB {
	shards := make([]*Shard, count)

	for i := 0; i < count; i++ {
		shards[i] = &Shard{data: make(map[string]string)}
	}

	return &DB{shardMap: shards}
}

func (d *DB) Get(key string) (string, bool) {
	shard := d.getShard(key)

	shard.RLock()
	defer shard.RUnlock()

	val, ok := shard.data[key]
	return val, ok
}

func (d *DB) Delete(key string) {
	shard := d.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	delete(shard.data, key)
}

func (d *DB) Set(key string, val string) {
	shard := d.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	shard.data[key] = val
}

func (d *DB) getShard(key string) *Shard {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[0])

	shardIndex := hash % len(d.shardMap)
	if shardIndex >= len(d.shardMap) {
		shardIndex %= len(d.shardMap)
	}

	return d.shardMap[shardIndex]
}
