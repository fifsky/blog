package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"app/pkg/errors"
	"app/proto/gen/types"

	"github.com/goapt/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoencoder = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

var errorEncoder = protojson.MarshalOptions{
	UseProtoNames: true,
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	switch vv := any(v).(type) {
	case *types.ErrorResponse:
		buf, err := errorEncoder.Marshal(vv)
		if err != nil {
			return fmt.Errorf("encode error protojson: %w", err)
		}
		_, err = w.Write(buf)
		return err
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
	resp := types.ErrorResponse_builder{}.Build()
	resp.SetCode(err.GetReason())
	resp.SetMessage(err.GetMessage())
	details := err.GetMetadata()
	if len(err.GetMetadata()) > 0 {
		details = err.GetMetadata()
	}

	if unErr := errors.Unwrap(err); unErr != nil {
		if details == nil {
			details = make(map[string]string)
		}
		details["cause"] = unErr.Error()
	}
	if len(details) > 0 {
		resp.SetDetails(details)
	}

	if e := encode(w, int(err.GetCode()), resp); e != nil {
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
