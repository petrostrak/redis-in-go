package main

import "sync"

type KV struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

func (kv *KV) Set(key, value []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[string(key)] = []byte(value)
	return nil
}

func (kv *KV) Get(key string) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	value, ok := kv.data[string(key)]
	return value, ok
}
