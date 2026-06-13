package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"app/pkg/errors"

	"github.com/goapt/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ErrorResponse struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

var protoencoder = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	switch vv := any(v).(type) {
	case *ErrorResponse:
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		return nil
	case proto.Message:
		buf, err := protoencoder.Marshal(vv)
		if err != nil {
			return fmt.Errorf("encode protojson: %w", err)
		}
		_, err = w.Write(buf)
		return err
	default:
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		return nil
	}
}

func Success(w http.ResponseWriter, data any) {
	err := encode(w, http.StatusOK, data)
	if err != nil {
		logger.Error("response error", slog.String("err", err.Error()))
	}
}

func Fail(w http.ResponseWriter, err *errors.Error) {
	resp := &ErrorResponse{
		Code:    err.Reason,
		Message: err.Message,
	}
	if len(err.Metadata) > 0 {
		resp.Details = err.Metadata
	}

	if unErr := errors.Unwrap(err); unErr != nil {
		if resp.Details == nil {
			resp.Details = make(map[string]string)
		}
		resp.Details["cause"] = unErr.Error()
	}

	if e := encode(w, int(err.Code), resp); e != nil {
		logger.Error("response error", slog.String("err", e.Error()))
	}
}

func FailPlain(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func Upload(w http.ResponseWriter, v any) {
	if err := encode(w, http.StatusOK, v); err != nil {
		logger.Error("response error", slog.String("err", err.Error()))
	}
}
