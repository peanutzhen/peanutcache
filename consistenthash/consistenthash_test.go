// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consistenthash

import (
	"hash/crc32"
	"log"
	"sort"
	"testing"
)

func TestConsistency_Register(t *testing.T) {
	c := New(2, nil)
	c.Register("peer1", "peer2")
	// 测试虚拟节点个数是否正确
	// Expect: replicas*len(peersName)
	if len(c.ring) != 4 {
		t.Errorf("Actual: %d\tExpect: %d\n", len(c.ring), 4)
	}
	// 测试哈希值是否正确
	hashValue := int(crc32.ChecksumIEEE([]byte("1peer1")))
	idx := sort.SearchInts(c.ring, hashValue)
	if c.ring[idx] != hashValue {
		t.Errorf("Actual: %d\tExpect: %d\n", c.ring[idx], hashValue)
	}
}

func TestConsistency_GetPeer(t *testing.T) {
	c := New(1, nil)
	c.Register("peer1", "peer2")
	key := "Tom"
	keyHashValue := int(crc32.ChecksumIEEE([]byte(key)))
	log.Printf("key hash = %d\n", keyHashValue)
	for _, v := range c.ring {
		log.Printf("%d -> %s\n", v, c.hashmap[v])
	}
	peer:=c.GetPeer(key)
	log.Printf("Go to search -> %s\n", peer)
}
