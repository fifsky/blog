package testutil

import (
	_ "github.com/go-sql-driver/mysql"
)

var TestDSN = "root:123456@tcp(127.0.0.1:3306)/"