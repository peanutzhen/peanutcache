// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package registry

//import (
//	"context"
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"testing"
//	"time"
//)

//func TestEtcdAdd(t *testing.T) {
//	cli, _ := clientv3.New(clientv3.Config{
//		Endpoints:   []string{"localhost:2379"},
//		DialTimeout: 5 * time.Second,
//	})
//	// 创建一个租约 配置5秒过期
//	resp, err := cli.Grant(context.Background(), 5)
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	err = EtcdAdd(cli, resp.ID, "test", "127.0.0.1:6324")
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//}
