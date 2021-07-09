package handler

import (
	"net/http"
	"testing"

	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestComment_Avatar(t *testing.T) {
	handler := NewCommon()
	req := test.NewRequest("/api/avatar", handler.Avatar)
	resp, err := req.JSON(``)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
}
