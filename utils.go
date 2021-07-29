// Copyright 2021 Peanutzhen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package peanutcache

import (
	"fmt"
	"runtime"
	"strings"
)

// 显示错误时运行堆栈
func trace(errorMessage string) string {
	var pcstack [32]uintptr
	n := runtime.Callers(3, pcstack[:])

	// Using Builder optimize speed.
	var str strings.Builder
	str.WriteString(errorMessage + "\nTraceback:")
	for _, pc := range pcstack[:n] {
		function := runtime.FuncForPC(pc)
		file, line := function.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
