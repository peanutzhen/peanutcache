// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func createTestSvr() *httptest.Server {
	mysql := map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	NewGroup("scores", 2<<10, RetrieverFunc(
		func(key string) ([]byte, error) {
			log.Println("[Mysql] search key", key)
			if v, ok := mysql[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	svr, _ := NewServer("localhost:9999")
	ts := httptest.NewServer(svr)
	svr.addr = ts.URL
	return ts
}

func TestServer_GetExistsKey(t *testing.T) {
	ts := createTestSvr()
	res, _ := http.Get(fmt.Sprintf("%s%sscores/Tom", ts.URL, defaultBasePath))
	body, _ := ioutil.ReadAll(res.Body)
	if !reflect.DeepEqual(string(body), "630") {
		t.Errorf("Tom %s(actual)/%s(ok)", string(body), "630")
	}
	res.Body.Close()
	ts.Close()
}

func TestServer_GetBadPath(t *testing.T) {
	ts := createTestSvr()
	res, _ := http.Get(fmt.Sprintf("%s%sscores/Tom", ts.URL, "/fakerBasePath/"))
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code should be 500.")
	}
	body, _ := ioutil.ReadAll(res.Body)
	t.Log("错误basePath查询Tom返回: " + string(body))
	res.Body.Close()
	ts.Close()
}

func TestServer_GetUnknownKey(t *testing.T) {
	ts := createTestSvr()
	res, _ := http.Get(fmt.Sprintf("%s%sscores/Unknown", ts.URL, defaultBasePath))
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code should be 500.")
	}
	body, _ := ioutil.ReadAll(res.Body)
	t.Log("正确规则查询不存在key: " + string(body))
	res.Body.Close()
}

func TestServer_GetUnknownGroup(t *testing.T) {
	ts := createTestSvr()
	res, _ := http.Get(fmt.Sprintf("%s%sUnknown/Tom", ts.URL, defaultBasePath))
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code should be 500.")
	}
	body, _ := ioutil.ReadAll(res.Body)
	t.Log("正确规则查询不存在group: " + string(body))
	res.Body.Close()
}

func TestServer_GetNoKey(t *testing.T) {
	ts := createTestSvr()
	res, _ := http.Get(fmt.Sprintf("%s%sUnknown", ts.URL, defaultBasePath))
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code should be 500.")
	}
	body, _ := ioutil.ReadAll(res.Body)
	t.Log("错误规则不填key返回: " + string(body))
	res.Body.Close()
}
