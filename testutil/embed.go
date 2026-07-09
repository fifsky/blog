package testutil

import (
	"bytes"
	"embed"
	"io"
	"path/filepath"
)

//go:embed testdata/schema.sql
var schemaBytes []byte

//go:embed testdata/fixtures
var fixturesFS embed.FS

//go:embed testdata/test.png
var testImageBytes []byte

// Schema 返回嵌入的 schema.sql 作为 io.Reader
func Schema() io.Reader {
	return bytes.NewReader(schemaBytes)
}

// Fixture 返回指定 fixture 文件的嵌入内容
func Fixture(file string) []byte {
	data, err := fixturesFS.ReadFile(filepath.Join("testdata", "fixtures", file+".yml"))
	if err != nil {
		panic("testutil: fixture not found: " + file + ".yml: " + err.Error())
	}
	return data
}

// Fixtures 返回多个 fixture 文件的嵌入内容（map[文件名]内容）
func Fixtures(files ...string) map[string][]byte {
	result := make(map[string][]byte, len(files))
	for _, file := range files {
		result[file+".yml"] = Fixture(file)
	}
	return result
}

// AllFixtures 返回所有嵌入的 fixture 文件
func AllFixtures() map[string][]byte {
	entries, err := fixturesFS.ReadDir("testdata/fixtures")
	if err != nil {
		panic("testutil: read fixtures dir error: " + err.Error())
	}
	result := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := fixturesFS.ReadFile(filepath.Join("testdata", "fixtures", entry.Name()))
		if err != nil {
			panic("testutil: read fixture error: " + entry.Name() + ": " + err.Error())
		}
		result[entry.Name()] = data
	}
	return result
}

// Image 返回嵌入的 test.png 测试图片内容
func Image() []byte {
	return testImageBytes
}
