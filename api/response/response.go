package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

type ApiResponse struct {
	Code int    `json:"code"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg"`
}

func Success(w http.ResponseWriter, data any) {
	err := encode(w, http.StatusOK, &ApiResponse{
		Code: 200,
		Data: data,
		Msg:  "success",
	})
	if err != nil {
		slog.Error("response error", slog.String("err", err.Error()))
	}
}

func Fail(w http.ResponseWriter, code int, msg any) {
	var m string
	switch e := msg.(type) {
	case string:
		m = e
	case error:
		m = e.Error()
	default:
		m = fmt.Sprintf("%v", e)
	}

	err := encode(w, 400, &ApiResponse{
		Code: code,
		Msg:  m,
	})

	if err != nil {
		slog.Error("response error", slog.String("err", err.Error()))
	}
}

func Upload(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("response error", slog.String("err", err.Error()))
	}
}
