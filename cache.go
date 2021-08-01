// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

// cache 模块负责提供对lru模块的并发控制

import (
	"github.com/peanutzhen/peanutcache/lru"
	"sync"
)

// 这样设计可以进行cache和算法的分离，比如我现在实现了lfu缓存模块
// 只需替换cache成员即可
type cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	capacity int64 // 缓存最大容量
}

func newCache(capacity int64) *cache {
	return &cache{capacity: capacity}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.capacity, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	if c.lru == nil {
		return ByteView{}, false
	}
	// 注意：Get操作需要修改lru中的双向链表，需要使用互斥锁。
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}
	return ByteView{}, false
}
