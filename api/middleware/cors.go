package middleware

import (
	"net/http"
	"strings"

	"app/config"
)

type Cors = func(next http.Handler) http.Handler

func NewCors(conf *config.Config) Cors {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origins := []string{"http://fifsky.com", "http://www.fifsky.com", "https://fifsky.com", "https://www.fifsky.com"}

			if conf.Env == "local" {
				origins = []string{"*"}
			}

			// 设置CORS头部
			w.Header().Set("Access-Control-Allow-Origin", getOrigin(r.Header.Get("Origin"), origins))
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Access-Token,Access-Control-Allow-Origin,x-requested-with")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "43200") // 12小时

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getOrigin 检查并返回匹配的源地址
func getOrigin(requestedOrigin string, allowedOrigins []string) string {
	for _, origin := range allowedOrigins {
		if origin == "*" {
			return requestedOrigin
		}
		if strings.ToLower(origin) == strings.ToLower(requestedOrigin) {
			return requestedOrigin
		}
	}
	return "" // 如果没有匹配的源，则返回空字符串
}
