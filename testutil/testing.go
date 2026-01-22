package testutil

import (
	"path/filepath"
	"runtime"
	"strings"
)

func TestDataPath(path ...string) string {
	_, file, _, _ := runtime.Caller(0)
	paths := []string{filepath.Dir(filepath.Dir(file)), "testdata"}
	paths = append(paths, path...)

	return filepath.Join(paths...)
}

func Schema() string {
	return filepath.Join(TestDataPath(), "schema-postgre.sql")
}

func Fixture(file string) string {
	return filepath.Join(FixturePath(), strings.TrimSuffix(file, ".yml")+".yml")
}

func FixturePath() string {
	s := Schema()
	return filepath.Join(filepath.Dir(s), "fixtures")
}

func Fixtures(files ...string) []string {
	for k, file := range files {
		files[k] = Fixture(file)
	}
	return files
}

func FixturesWithPath(path string, files ...string) []string {
	for k, file := range files {
		files[k] = filepath.Join(TestDataPath(path), strings.TrimSuffix(file, ".yml")+".yml")
	}
	return files
}
