package store

import (
	"log"
	"sync"
)

type entry struct {
	value  string
	bitmap string
}

var signatureStore = struct {
	sync.RWMutex
	m map[string]entry
}{m: make(map[string]entry)}

func Put(key string, value, bitmap string) {
	log.Printf("put: %s, %v", key, value)
	signatureStore.Lock()
	signatureStore.m[key] = entry{
		value:  value,
		bitmap: bitmap,
	}
	signatureStore.Unlock()
}

func Get(key string) (string, string, bool) {
	val, ok := signatureStore.m[key]
	return val.value, val.bitmap, ok
}
