package openapi

import (
	"context"
	"testing"

	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestSetting_Get(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		svc := NewSetting(store.New(db))
		resp, err := svc.Get(context.Background(), &emptypb.Empty{})
		if err != nil || len(resp.Kv) == 0 {
			t.Fatalf("unexpected err=%v kv=%v", err, resp.Kv)
		}
	})
}

func TestSetting_GetDoesNotExposeAIToken(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		_, err := db.ExecContext(context.Background(), "insert into options (option_key, option_value) values (?, ?)", "ai_token", "secret-token")
		require.NoError(t, err)

		svc := NewSetting(store.New(db))
		resp, err := svc.Get(context.Background(), &emptypb.Empty{})
		require.NoError(t, err)
		assert.NotContains(t, resp.Kv, "ai_token")
	})
}
