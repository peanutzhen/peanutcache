// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

// peers 模块

// Picker 定义了获取分布式节点的能力
type Picker interface {
	Pick(key string) (Fetcher, bool)
}

// Fetcher 定义了从远端获取缓存的能力
// 所以每个Peer应实现这个接口
type Fetcher interface {
	Fetch(group string, key string) ([]byte, error)
}
