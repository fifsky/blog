package gravatar

import (
	"strings"
	"testing"
)

func TestAvatarURL(t *testing.T) {
	tests := []struct {
		name  string
		email string
		size  int
		want  string
	}{
		{
			name:  "标准邮箱",
			email: "My.Email@example.com ", // 大小写+空白，应规范化
			size:  80,
			want:  GravatarProxy + "/5286fb0787d7a84ac13c19f44ef86bc7?s=80",
		},
		{
			name:  "默认尺寸",
			email: "test@example.com",
			size:  0,
			want:  GravatarProxy + "/55502f40dc8b7c769880b10874abc9d0?s=80",
		},
		{
			name:  "空邮箱返回占位头像",
			email: "",
			size:  40,
			want:  GravatarProxy + "/d41d8cd98f00b204e9800998ecf8427e?s=40",
		},
		{
			name:  "负尺寸使用默认值",
			email: "a@b.com",
			size:  -1,
			want:  GravatarProxy + "/357a20e8c56e69d6f9734d23ef9517e8?s=80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AvatarURL(tt.email, tt.size)
			if got != tt.want {
				t.Errorf("AvatarURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAvatarURL_Format 验证生成的 URL 格式符合代理地址规范
func TestAvatarURL_Format(t *testing.T) {
	got := AvatarURL("someone@example.com", 120)
	if !strings.HasPrefix(got, GravatarProxy+"/") {
		t.Errorf("URL 应以代理前缀开头，实际: %s", got)
	}
	if !strings.HasSuffix(got, "?s=120") {
		t.Errorf("URL 应以 ?s=120 结尾，实际: %s", got)
	}
}
