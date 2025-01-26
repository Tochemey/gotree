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
	"hash/fnv"
	"sync"
)

// maps is a concurrent map with sharding for scalability
type maps struct {
	shards    []*sync.Map
	numShards uint64
}

// newMaps creates an instance of maps
func newMaps(totalShards uint64) *maps {
	shards := make([]*sync.Map, totalShards)
	for i := range shards {
		shards[i] = &sync.Map{}
	}
	return &maps{
		shards:    shards,
		numShards: totalShards,
	}
}

// getShard returns the given shard for a given key
func (s *maps) getShard(key string) *sync.Map {
	hash := fnv64(key) % s.numShards
	return s.shards[hash]
}

// Load returns the value of a given key
func (s *maps) Load(key string) (any, bool) {
	shard := s.getShard(key)
	return shard.Load(key)
}

// Store adds a key/value pair to the sharded map
func (s *maps) Store(key string, value any) {
	shard := s.getShard(key)
	shard.Store(key, value)
}

// Delete removes a given key from the sharded map
func (s *maps) Delete(key string) {
	shard := s.getShard(key)
	shard.Delete(key)
}

// Range given a function iterate over the sharded ma[
func (s *maps) Range(f func(key, value any) bool) {
	for i := 0; i < int(s.numShards); i++ {
		shard := s.shards[i]
		shard.Range(f)
	}
}

// Reset resets the sharded map
func (s *maps) Reset() {
	// Reset each shard's map
	for i := 0; i < int(s.numShards); i++ {
		shard := s.shards[i]
		shard.Range(func(key, _ any) bool {
			shard.Delete(key) // Clear the entry
			return true
		})
	}
}

func fnv64(key string) uint64 {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(key))
	return hash.Sum64()
}
