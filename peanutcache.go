// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"log"
	"sync"
)

// peanutcache 模块提供比cache模块更高一层抽象的能力
// 换句话说，实现了填充缓存/命名划分缓存的能力

var (
	mu     sync.RWMutex // 管理读写groups并发控制
	groups = make(map[string]*Group)
)

// Retriever 要求对象实现从数据源获取数据的能力
type Retriever interface {
	retrieve(string) ([]byte, error)
}

type RetrieverFunc func(key string) ([]byte, error)

// RetrieverFunc 通过实现retrieve方法，使得任意匿名函数func
// 通过被RetrieverFunc(func)类型强制转换后，实现了 Retriever 接口的能力
func (f RetrieverFunc) retrieve(key string) ([]byte, error) {
	return f(key)
}

// Group 提供命名管理缓存/填充缓存的能力
type Group struct {
	name      string
	cache     *cache
	retriever Retriever
}

// NewGroup 创建一个新的缓存空间
func NewGroup(name string, maxBytes int64, retriever Retriever) *Group {
	if retriever == nil {
		panic("Group retriever must be existed!")
	}
	g := &Group{
		name:      name,
		cache:     newCache(maxBytes),
		retriever: retriever,
	}
	mu.Lock()
	groups[name] = g
	mu.Unlock()
	return g
}

// GetGroup 获取对应命名空间的缓存
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key required")
	}
	if value, ok := g.cache.get(key); ok {
		log.Println("cache hit")
		return value, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.retriever.retrieve(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.cache.add(key, value)
}