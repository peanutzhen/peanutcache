// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consistenthash

// consistenthash 模块负责实现一致性哈希
// 用于确定key与peer之间的映射

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashFunc 定义哈希函数输入输出
type HashFunc func(data []byte) uint32

// Consistency 维护peer与其hash值的关联
type Consistency struct {
	hash     HashFunc       // 哈希函数依赖
	replicas int            // 虚拟节点个数(防止数据倾斜)
	ring     []int          // uint32哈希环
	hashmap  map[int]string // hashValue -> peerName
}

// Register 将各个peer注册到哈希环上
func (c *Consistency) Register(peersName ...string) {
	for _, peerName := range peersName {
		for i := 0; i < c.replicas; i++ {
			hashValue := int(c.hash([]byte(strconv.Itoa(i)+peerName)))
			c.ring = append(c.ring, hashValue)
			c.hashmap[hashValue] = peerName
		}
	}
	sort.Ints(c.ring)
}

// GetPeer 计算key应缓存到的peer
func (c *Consistency) GetPeer(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hashValue := int(c.hash([]byte(key)))
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hashValue
	})
	return c.hashmap[c.ring[idx%len(c.ring)]]
}

func New(replicas int, fn HashFunc) *Consistency {
	c := &Consistency{
		replicas: replicas,
		hash:     fn,
		hashmap:  make(map[int]string),
	}
	if c.hash == nil {
		c.hash = crc32.ChecksumIEEE
	}
	return c
}
