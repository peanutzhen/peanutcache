// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// server 模块为peanutcache提供http通信能力
// 这样部署在其他机器上的cache可以通过http访问获取缓存
// 至于找哪台主机 那是一致性哈希的工作了

const defaultBasePath = "/_pcache/"

type Server struct {
	addr     string
	basePath string
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		basePath: defaultBasePath,
	}
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
		panic("Key required")
	}

	groupName, key := pathsName[0], pathsName[1]
	log.Printf("[peanutcache %s] GET - (%s)/(%s)", s.addr, groupName, key)

	g := GetGroup(groupName)
	if g == nil {
		panic("group not found")
	}

	view, err := g.Get(key)
	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(view.ByteSlice())
	if err != nil {
		panic("ByteView write to response failed")
	}
}