package admin

import (
	"context"
	"testing"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/require"
)

func TestAdminSetting_Update(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		svc := NewSetting(store.New(db))

		resp, err := svc.Update(context.Background(), &adminv1.AdminSetting{
			SiteName: "abc",
			SiteDesc: "def",
		})
		require.NoError(t, err)
		require.Equal(t, "abc", resp.SiteName)
		require.Equal(t, "def", resp.SiteDesc)
	})
}
