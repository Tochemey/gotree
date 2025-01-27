/*
 * MIT License
 *
 * Copyright (c) 2025 Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gotree

import (
	"crypto/sha1"
	"hash/fnv"
	"sync"
)

// Shard defines a Shard
type Shard struct {
	sync.RWMutex
	m map[string]any
}

// ShardedMap defines a concurrent map with sharding for
// scalability
type ShardedMap []*Shard

// NewShardedMap creates an instance of ShardedMap
func NewShardedMap(shardsCount uint64) ShardedMap {
	shards := make([]*Shard, shardsCount)
	for i := range shardsCount {
		shards[i] = &Shard{
			m: make(map[string]any),
		}
	}
	return shards
}

// Load returns the value of a given key
func (s ShardedMap) Load(key string) (any, bool) {
	shard := s.getShard(key)
	shard.RLock()
	val, ok := shard.m[key]
	shard.RUnlock()
	return val, ok
}

// Store adds a key/value pair to the sharded map
func (s ShardedMap) Store(key string, value any) {
	shard := s.getShard(key)
	shard.Lock()
	shard.m[key] = value
	shard.Unlock()
}

// Delete removes a given key from the sharded map
func (s ShardedMap) Delete(key string) {
	shard := s.getShard(key)
	shard.Lock()
	delete(shard.m, key)
	shard.Unlock()
}

// Range given a function iterate over the sharded map
func (s ShardedMap) Range(f func(key, value any) bool) {
	for i := 0; i < len(s); i++ {
		shard := s[i]
		shard.RLock()
		for k, v := range shard.m {
			f(k, v)
		}
		shard.RUnlock()
	}
}

// Reset resets the sharded map
func (s ShardedMap) Reset() {
	// Reset each Shard's map
	for i := 0; i < len(s); i++ {
		shard := s[i]
		shard.Lock()
		shard.m = make(map[string]any)
		shard.Unlock()
	}
}

// getShard returns the given Shard for a given key
func (s ShardedMap) getShard(key string) *Shard {
	hash := fnv64(key) % uint64(len(s))
	return s[int(hash)]
}

// shardIndex returns the Shard index
func (s ShardedMap) shardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(s)
}

func fnv64(key string) uint64 {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(key))
	return hash.Sum64()
}
