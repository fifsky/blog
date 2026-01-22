package testutil

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

var TestDSN = "postgresql://postgres:123456@localhost:5432/?timezone=Asia%2FShanghai"
