// Package gravatar 提供基于邮箱的 Gravatar 头像代理地址生成。
package gravatar

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

// GravatarProxy 是 Gravatar 头像代理地址前缀
const GravatarProxy = "https://seccdn.libravatar.org/gravatarproxy"

// DefaultSize 默认头像尺寸
const DefaultSize = 80

// AvatarURL 根据邮箱生成 Gravatar 头像代理地址。
// 邮箱为空时使用空字符串的哈希，返回一个默认占位头像。
func AvatarURL(email string, size int) string {
	if size <= 0 {
		size = DefaultSize
	}
	// Gravatar 规范：先去除首尾空白，再转小写，最后取 md5
	sum := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	hash := hex.EncodeToString(sum[:])
	return GravatarProxy + "/" + hash + "?s=" + strconv.Itoa(size)
}
