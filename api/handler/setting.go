package handler

import (
	"net/http"

	"app/response"
	"app/store"
)

type Setting struct {
	store *store.Store
}

func NewSetting(s *store.Store) *Setting {
	return &Setting{store: s}
}

func (s *Setting) Get(w http.ResponseWriter, r *http.Request) {
	m, err := s.store.GetOptions(r.Context())
	if err != nil {
		response.Fail(w, 202, err)
		return
	}

	var resp OptionsResponse = m
	response.Success(w, resp)
}

func (s *Setting) Post(w http.ResponseWriter, r *http.Request) {
	kv, err := decode[map[string]string](r)
	if err != nil {
		response.Fail(w, 202, err)
		return
	}

	m, err := s.store.UpdateOptions(r.Context(), kv)
	if err != nil {
		response.Fail(w, 203, err)
		return
	}
	var resp OptionsResponse = m
	response.Success(w, resp)
}
