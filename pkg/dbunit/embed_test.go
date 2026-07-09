package dbunit

import (
	"bytes"
	"embed"
	"io"
)

//go:embed testdata/schema.sql
var testSchemaBytes []byte

//go:embed testdata/fixtures
var testFixturesFS embed.FS

//go:embed testdata/custom
var testCustomFS embed.FS

// testSchemaReader 返回 dbunit 自身测试用 schema 的 io.Reader
func testSchemaReader() io.Reader {
	return bytes.NewReader(testSchemaBytes)
}

// testFixturesMap 从嵌入的 fixtures 目录构建 fixture 内容映射
func testFixturesMap(names ...string) map[string][]byte {
	result := make(map[string][]byte, len(names))
	for _, n := range names {
		data, err := testFixturesFS.ReadFile("testdata/fixtures/" + n + ".yml")
		if err != nil {
			panic("dbunit test: fixture not found: " + n + ".yml: " + err.Error())
		}
		result[n+".yml"] = data
	}
	return result
}

// testCustomFixtures 从嵌入的 custom 目录构建 fixture 内容映射
func testCustomFixtures() map[string][]byte {
	result := make(map[string][]byte)
	entries, err := testCustomFS.ReadDir("testdata/custom")
	if err != nil {
		panic("dbunit test: read custom dir error: " + err.Error())
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := testCustomFS.ReadFile("testdata/custom/" + entry.Name())
		if err != nil {
			panic("dbunit test: read custom fixture error: " + entry.Name() + ": " + err.Error())
		}
		result[entry.Name()] = data
	}
	return result
}
