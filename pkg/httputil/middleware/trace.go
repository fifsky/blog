// Copyright (c) 2015-2021 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Trace() Middleware {
	return func(tr http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(tr)
	}
}
