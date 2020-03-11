/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringWeb_test

import (
	"net/http"
	"testing"

	"github.com/go-spring/go-spring-web/spring-web"
)

// methods
var methods = map[uint32]string{
	SpringWeb.MethodGet:     http.MethodGet,
	SpringWeb.MethodHead:    http.MethodHead,
	SpringWeb.MethodPost:    http.MethodPost,
	SpringWeb.MethodPut:     http.MethodPut,
	SpringWeb.MethodPatch:   http.MethodPatch,
	SpringWeb.MethodDelete:  http.MethodDelete,
	SpringWeb.MethodConnect: http.MethodConnect,
	SpringWeb.MethodOptions: http.MethodOptions,
	SpringWeb.MethodTrace:   http.MethodTrace,
}

// cacheMethods
var cacheMethods = map[uint32][]string{
	SpringWeb.MethodGet:     {http.MethodGet},
	SpringWeb.MethodHead:    {http.MethodHead},
	SpringWeb.MethodPost:    {http.MethodPost},
	SpringWeb.MethodPut:     {http.MethodPut},
	SpringWeb.MethodPatch:   {http.MethodPatch},
	SpringWeb.MethodDelete:  {http.MethodDelete},
	SpringWeb.MethodConnect: {http.MethodConnect},
	SpringWeb.MethodOptions: {http.MethodOptions},
	SpringWeb.MethodTrace:   {http.MethodTrace},
	SpringWeb.MethodGetPost: {http.MethodGet, http.MethodPost},
	SpringWeb.MethodAny:     {http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace},
}

func getMethod(method uint32) []string {
	var r []string
	for k, v := range methods {
		if method&k == k {
			r = append(r, v)
		}
	}
	return r
}

func getMethodViaCache(method uint32) []string {
	if r, ok := cacheMethods[method]; ok {
		return r
	}
	return getMethod(method)
}

func BenchmarkGetMethod(b *testing.B) {
	// 事实证明通过缓存不一定能提高效率

	b.Run("1", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet)
	})

	b.Run("cache-1", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet)
	})

	b.Run("2", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead)
	})

	b.Run("cache-2", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead)
	})

	b.Run("3", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost)
	})

	b.Run("cache-3", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost)
	})

	b.Run("4", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut)
	})

	b.Run("cache-4", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut)
	})

	b.Run("5", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch)
	})

	b.Run("cache-5", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch)
	})

	b.Run("6", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete)
	})

	b.Run("cache-6", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete)
	})

	b.Run("7", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect)
	})

	b.Run("cache-7", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect)
	})

	b.Run("8", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect | SpringWeb.MethodOptions)
	})

	b.Run("cache-8", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect | SpringWeb.MethodOptions)
	})

	b.Run("9", func(b *testing.B) {
		getMethod(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect | SpringWeb.MethodOptions | SpringWeb.MethodTrace)
	})

	b.Run("cache-9", func(b *testing.B) {
		getMethodViaCache(SpringWeb.MethodGet | SpringWeb.MethodHead | SpringWeb.MethodPost | SpringWeb.MethodPut | SpringWeb.MethodPatch | SpringWeb.MethodDelete | SpringWeb.MethodConnect | SpringWeb.MethodOptions | SpringWeb.MethodTrace)
	})
}