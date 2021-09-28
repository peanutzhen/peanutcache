# PeanutCache

![](https://img.shields.io/badge/license-MIT-blue)![](https://img.shields.io/github/stars/peanutzhen/peanutcache?style=plastic)

The [**gRPC**](https://github.com/grpc/grpc-go) implementation of [**groupcache**](https://github.com/golang/groupcache): A high performance, open source, using RPC framework that  communicated with each cache node. Cache service can register to [**etcd**](https://github.com/etcd-io/etcd), and each cache client can dicovery the service list by etcd.For more information see the [groupcache](https://github.com/golang/groupcache), or [**geecache**](https://geektutu.com/post/geecache.html).

## Prerequisites

- **Golang** 1.16 or later
- **Etcd** v3.4.0 or later
- **gRPC-go** v1.38.0 or later
- **protobuf** v1.26.0 or later

## TodoList

欢迎大家Pull Request，可随时联系作者。

1. 将一致性哈希从`Server`抽象出来，作为单独的一个`Proxy`层。避免在每个节点自己做一致性哈希，这样存在哈希环不一致的情况。
2. 增加缓存持久化的能力。
3. 改进`LRU cache`，使其具备`TTL`的能力，以及改进锁的粒度，提高并发度。

## Installation

With [Go module]() support (Go 1.11+), simply add the following import

```go
import "github.com/peanutzhen/peanutcache"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, to install the `peanutcache` package, run the following command:

```bash
$ go get -u github.com/peanutzhen/peanutcache
```

## Usage

Here, give a example to use it `example.go`:

```go
// example.go file
// 运行前，你需要在本地启动Etcd实例，作为服务中心。

package main

import (
	"fmt"
	"github.com/peanutzhen/peanutcache"
	"log"
	"sync"
)

func main() {
	// 模拟MySQL数据库 用于peanutcache从数据源获取值
	var mysql = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	// 新建cache实例
	group := peanutcache.NewGroup("scores", 2<<10, peanutcache.RetrieverFunc(
		func(key string) ([]byte, error) {
			log.Println("[Mysql] search key", key)
			if v, ok := mysql[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	// New一个服务实例
	var addr string = "localhost:9999"
	svr, err := peanutcache.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	// 设置同伴节点IP(包括自己)
  // todo: 这里的peer地址从etcd获取(服务发现)
	svr.SetPeers(addr)
	// 将服务与cache绑定 因为cache和server是解耦合的
	group.RegisterSvr(svr)
	log.Println("peanutcache is running at", addr)
	// 启动服务(注册服务至etcd/计算一致性哈希...)
	go func() {
		// Start将不会return 除非服务stop或者抛出error
		err = svr.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 发出几个Get请求
	var wg sync.WaitGroup
	wg.Add(4)
	go GetTomScore(group, &wg)
	go GetTomScore(group, &wg)
	go GetTomScore(group, &wg)
	go GetTomScore(group, &wg)
	wg.Wait()
}

func GetTomScore(group *peanutcache.Group, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("get Tom...")
	view, err := group.Get("Tom")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(view.String())
}
```
Before `go run`, you should run `etcd` local directly(without any spcified parameter) and then execute `go run example.go`, you will get follows:

```console
$ go run peanutcache_usage.go
2021/08/01 18:09:17 peanutcache is running at localhost:9999
2021/08/01 18:09:17 get Tom...
2021/08/01 18:09:17 get Tom...
2021/08/01 18:09:17 get Tom...
2021/08/01 18:09:17 get Tom...
2021/08/01 18:09:17 ooh! pick myself, I am localhost:9999
2021/08/01 18:09:17 [Mysql] search key Tom
630
630
630
630
```

