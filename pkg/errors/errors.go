package errors

import (
	"errors"
	"fmt"
	"maps"

	httpstatus "app/pkg/status"
	"app/proto/gen/types"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is a status error.
type Error struct {
	types.Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.GetCode(), e.GetReason(), e.GetMessage(), e.GetMetadata(), e.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.GetCode() == e.GetCode() && se.GetReason() == e.GetReason()
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.SetMetadata(md)
	return err
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.GetCode())), e.GetMessage()).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.GetReason(),
			Metadata: e.GetMetadata(),
		})
	return s
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	err := &Error{}
	err.SetCode(int32(code))
	err.SetReason(reason)
	err.SetMessage(message)
	return err
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...any) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...any) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:mnd
	}
	return int(FromError(err).GetCode())
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).GetReason()
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.GetMetadata()))
	maps.Copy(metadata, err.GetMetadata())
	ret := New(int(err.GetCode()), err.GetReason(), err.GetMessage())
	ret.cause = err.cause
	ret.SetMetadata(metadata)
	return ret
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, UnknownReason, err.Error())
	}
	ret := New(
		httpstatus.FromGRPCCode(gs.Code()),
		UnknownReason,
		gs.Message(),
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.SetReason(d.Reason)
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}
