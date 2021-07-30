// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// client 模块实现peanutcache访问其他远程节点
// 从而获取缓存的能力

type Client struct {
	reqURL string
}

// Fetch 从remote peer获取对应缓存值
func (c *Client) Fetch(group string, key string) ([]byte, error) {
	// 构造请求url
	u := fmt.Sprintf(
		"%s%s/%s",
		c.reqURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("peer Statuscode: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed, %v", err)
	}
	return body, nil
}

func NewClient(reqURL string) *Client {
	return &Client{reqURL: reqURL}
}

var _ Fetcher = (*Client)(nil)
