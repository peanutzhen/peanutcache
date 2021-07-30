// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/peanutzhen/peanutcache/consistenthash"
)

// server 模块为peanutcache之间提供http通信能力
// 这样部署在其他机器上的cache可以通过http访问获取缓存
// 至于找哪台主机 那是一致性哈希的工作了
// 注意: peer间通信采用http协议

const (
	defaultBasePath = "/_pcache/"
	defaultReplicas = 50
)

// Server 和 Group 是解耦合的 所以Server要自己实现并发控制
type Server struct {
	addr     string // format: ip:port
	basePath string

	mu       sync.Mutex
	consHash *consistenthash.Consistency
	clients  map[string]*Client
}

func NewServer(addr string) (*Server, error) {
	if len(strings.Split(addr, ":")) != 2 {
		return nil, fmt.Errorf("server addr format-> ip:port")
	}
	return &Server{
		addr:     addr,
		basePath: defaultBasePath,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// server 错误恢复
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Printf("%s\n\n", trace(message))
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(message))
			if err != nil {
				panic("Error message write to response failed")
			}
		}
	}()

	// 路由规则必须如下所示
	// /basePath/groupName/key
	if !strings.HasPrefix(req.URL.Path, s.basePath) {
		panic("Unexpected request path: " + req.URL.Path)
	}
	// 跳过/basePath/ 即groupName/key
	pathsName := strings.SplitN(req.URL.Path[len(s.basePath):], "/", 2)
	if len(pathsName) != 2 {
		// TODO: 这里单元测试时file line不准确 原因未知
		panic("Key required")
	}

	groupName, key := pathsName[0], pathsName[1]
	log.Printf("[peanutcache %s] GET - (%s)/(%s)", s.addr, groupName, key)

	g := GetGroup(groupName)
	if g == nil {
		// TODO: 这里单元测试时file line不准确 原因未知
		panic("group not found")
	}

	view, err := g.Get(key)
	if err != nil {
		// TODO: 这里单元测试时file line不准确 原因未知
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(view.ByteSlice())
	if err != nil {
		panic("ByteView write to response failed")
	}
}

// SetPeers 将各个远端主机IP配置到Server里
// 这样Server就可以Pick他们了
// 注意: 此操作是*覆写*操作！
// 注意: peersIP必须满足 http://x.x.x.x:port的格式
func (s *Server) SetPeers(peersURL ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consHash = consistenthash.New(defaultReplicas, nil)
	s.consHash.Register(peersURL...)
	s.clients = make(map[string]*Client)
	for _, peerURL := range peersURL {
		if !validPeerURL(peerURL) {
			panic(fmt.Sprintf("[peer %s] using not a http protocol or containing a path.", peerURL))
		}
		s.clients[peerURL] = NewClient(peerURL + defaultBasePath)
	}
}

// Pick 根据一致性哈希选举出key应存放在的cache
// return false 代表从本地获取cache
func (s *Server) Pick(key string) (Fetcher, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	peerURL := s.consHash.GetPeer(key)
	u, err := url.Parse(peerURL)
	if err != nil {
		return nil, false
	}
	// Pick itself
	if u.Host == s.addr {
		log.Printf("ooh! pick myself, I am %s\n", s.addr)
		return nil, false
	}
	log.Printf("Pick remote peer: %s\n", peerURL)
	return s.clients[peerURL], true
}

var _ Picker = (*Server)(nil)
