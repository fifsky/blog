package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	apiv1 "app/proto/gen/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoencoder = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: true,
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	switch vv := any(v).(type) {
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
		slog.Error("response error", slog.String("err", err.Error()))
	}
}

func Fail(w http.ResponseWriter, code int32, msg any) {
	var m string
	switch e := msg.(type) {
	case string:
		m = e
	case error:
		m = e.Error()
	default:
		m = fmt.Sprintf("%v", e)
	}

	err := encode(w, 400, &apiv1.ErorResponse{
		Code: code,
		Msg:  m,
	})

	if err != nil {
		slog.Error("response error", slog.String("err", err.Error()))
	}
}

func FailPlain(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func Upload(w http.ResponseWriter, v any) {
	if err := encode(w, http.StatusOK, v); err != nil {
		slog.Error("response error", slog.String("err", err.Error()))
	}
}
