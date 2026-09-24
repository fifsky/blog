package middleware

import (
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// CorsConfig 表示 CORS 中间件的全部可配置项。
type CorsConfig struct {
	// AllowAllOrigins 允许任意来源，等价于 AllowOrigins 中包含 "*"
	AllowAllOrigins bool

	// AllowOrigins 允许跨域的来源列表，支持 "*" 以及 AllowWildcard 开启后的通配写法
	AllowOrigins []string

	// AllowOriginFunc 自定义来源校验函数，返回 true 表示放行
	AllowOriginFunc func(origin string) bool

	// AllowMethods 允许跨域使用的请求方法
	AllowMethods []string

	// AllowHeaders 允许跨域携带的请求头
	AllowHeaders []string

	// AllowCredentials 是否允许携带凭证（Cookie、Authorization 等）
	AllowCredentials bool

	// ExposeHeaders 允许前端 JS 读取的响应头
	ExposeHeaders []string

	// MaxAge 预检结果的缓存时长，精度为秒
	MaxAge time.Duration

	// AllowWildcard 允许来源中使用通配符，如 https://*.fifsky.com、http://localhost:*
	AllowWildcard bool

	// AllowBrowserExtensions 允许浏览器扩展来源（如 chrome-extension://）
	AllowBrowserExtensions bool

	// AllowWebSockets 允许 ws/wss 协议来源
	AllowWebSockets bool

	// AllowFiles 允许 file:// 协议来源（存在安全风险，确认需要时再开启）
	AllowFiles bool
}

// AddAllowMethods 追加允许的方法
func (c *CorsConfig) AddAllowMethods(methods ...string) {
	c.AllowMethods = append(c.AllowMethods, methods...)
}

// AddAllowHeaders 追加允许的请求头
func (c *CorsConfig) AddAllowHeaders(headers ...string) {
	c.AllowHeaders = append(c.AllowHeaders, headers...)
}

// AddExposeHeaders 追加暴露给前端的响应头
func (c *CorsConfig) AddExposeHeaders(headers ...string) {
	c.ExposeHeaders = append(c.ExposeHeaders, headers...)
}

// CorsOption 用于覆盖默认的 CORS 配置
type CorsOption func(c *CorsConfig)

// CorsOrigins 设置允许的来源列表
func CorsOrigins(o []string) CorsOption {
	return func(c *CorsConfig) {
		c.AllowOrigins = o
	}
}

// CorsMethods 设置允许的请求方法
func CorsMethods(o []string) CorsOption {
	return func(c *CorsConfig) {
		c.AllowMethods = o
	}
}

// CorsHeaders 追加允许的请求头
func CorsHeaders(o []string) CorsOption {
	return func(c *CorsConfig) {
		c.AllowHeaders = append(c.AllowHeaders, o...)
	}
}

// parseWildcardRules 将带通配符的来源解析为「前缀 + 后缀」的匹配规则。
// 例如 https://*.fifsky.com 解析为 ["https://", ".fifsky.com"]，
// http://localhost:* 解析为 ["http://localhost", "*"]（后缀为 * 表示只校验前缀）。
func (c CorsConfig) parseWildcardRules() [][]string {
	var wRules [][]string

	if !c.AllowWildcard {
		return wRules
	}

	for _, o := range c.AllowOrigins {
		if !strings.Contains(o, "*") {
			continue
		}

		if cnt := strings.Count(o, "*"); cnt > 1 {
			log.Fatalf("[cors] only one * is allowed in origin, got %d in %q", cnt, o)
		}

		i := strings.Index(o, "*")
		if i == 0 {
			wRules = append(wRules, []string{"*", o[1:]})
			continue
		}
		if i == len(o)-1 {
			wRules = append(wRules, []string{o[:i-1], "*"})
			continue
		}

		wRules = append(wRules, []string{o[:i], o[i+1:]})
	}

	return wRules
}

// corsHandler 持有解析后的 CORS 配置，避免每次请求重复计算响应头
type corsHandler struct {
	allowAllOrigins  bool
	allowCredentials bool
	allowOriginFunc  func(string) bool
	allowOrigins     []string
	normalHeaders    http.Header
	preflightHeaders http.Header
	wildcardOrigins  [][]string
}

func newCorsHandler(config CorsConfig) *corsHandler {
	for _, origin := range config.AllowOrigins {
		if origin == "*" {
			config.AllowAllOrigins = true
		}
	}

	return &corsHandler{
		allowOriginFunc:  config.AllowOriginFunc,
		allowAllOrigins:  config.AllowAllOrigins,
		allowCredentials: config.AllowCredentials,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalHeaders(config),
		preflightHeaders: generatePreflightHeaders(config),
		wildcardOrigins:  config.parseWildcardRules(),
	}
}

// validateWildcardOrigin 按通配规则校验来源
func (ch *corsHandler) validateWildcardOrigin(origin string) bool {
	for _, w := range ch.wildcardOrigins {
		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
			return true
		}
		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
			return true
		}
		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
			return true
		}
	}

	return false
}

// validateOrigin 校验来源是否在白名单内
func (ch *corsHandler) validateOrigin(origin string) bool {
	if ch.allowAllOrigins {
		return true
	}
	if slices.Contains(ch.allowOrigins, origin) {
		return true
	}
	if len(ch.wildcardOrigins) > 0 && ch.validateWildcardOrigin(origin) {
		return true
	}
	if ch.allowOriginFunc != nil {
		return ch.allowOriginFunc(origin)
	}
	return false
}

// handlePreflight 写入预检响应头
func (ch *corsHandler) handlePreflight(w http.ResponseWriter) {
	maps.Copy(w.Header(), ch.preflightHeaders)
}

// handleNormal 写入普通请求的响应头
func (ch *corsHandler) handleNormal(w http.ResponseWriter) {
	maps.Copy(w.Header(), ch.normalHeaders)
}

// NewCors 构建 CORS 中间件，默认放行博客站点域名与本地开发端口，可通过 CorsOption 覆盖。
//
// OPTIONS 请求统一在此收口：业务路由只注册了带方法的模式（如 "POST /blog/xxx"），
// 探测型 OPTIONS 落到 ServeMux 会因方法不匹配直接返回 405，因此 CORS 需要包裹住整个 mux，
// 由中间件尽早终结预检请求（非跨域 204 / 预检 204），不再需要单独注册 OPTIONS 路由。
func NewCors(opts ...CorsOption) func(next http.Handler) http.Handler {
	conf := CorsConfig{
		AllowOrigins: []string{
			"https://fifsky.com",
			"http://fifsky.com",
			"https://*.fifsky.com",
			"http://*.fifsky.com",
			"http://localhost:*",
			"http://127.0.0.1:*",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Accept-Encoding",
			"Authorization",
			"x-requested-with",
			"X-Client-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
	}

	for _, o := range opts {
		o(&conf)
	}

	ch := newCorsHandler(conf)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			host := r.Host

			// 非跨域请求（无 Origin，或 Origin 与 Host 同源）不携带任何 CORS 头
			if len(origin) == 0 || origin == "http://"+host || origin == "https://"+host {
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if !ch.validateOrigin(origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if r.Method == http.MethodOptions {
				ch.handlePreflight(w)
				if !ch.allowAllOrigins {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			ch.handleNormal(w)
			if !ch.allowAllOrigins {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// generateNormalHeaders 生成普通请求的 CORS 响应头
func generateNormalHeaders(c CorsConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposeHeaders) > 0 {
		exposeHeaders := convert(normalize(c.ExposeHeaders), http.CanonicalHeaderKey)
		headers.Set("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ","))
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Vary", "Origin")
	}
	return headers
}

// generatePreflightHeaders 生成预检请求的 CORS 响应头
func generatePreflightHeaders(c CorsConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.AllowMethods) > 0 {
		allowMethods := convert(normalize(c.AllowMethods), strings.ToUpper)
		value := strings.Join(allowMethods, ",")
		headers.Set("Access-Control-Allow-Methods", value)
	}
	if len(c.AllowHeaders) > 0 {
		allowHeaders := convert(normalize(c.AllowHeaders), http.CanonicalHeaderKey)
		value := strings.Join(allowHeaders, ",")
		headers.Set("Access-Control-Allow-Headers", value)
	}
	if c.MaxAge > time.Duration(0) {
		value := strconv.FormatInt(int64(c.MaxAge/time.Second), 10)
		headers.Set("Access-Control-Max-Age", value)
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
	}
	return headers
}

// normalize 去空格、转小写并去重
func normalize(values []string) []string {
	if values == nil {
		return nil
	}
	distinctMap := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)
		if _, seen := distinctMap[value]; !seen {
			normalized = append(normalized, value)
			distinctMap[value] = true
		}
	}
	return normalized
}

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

// converter 字符串转换函数，如 strings.ToUpper、http.CanonicalHeaderKey
type converter func(string) string
