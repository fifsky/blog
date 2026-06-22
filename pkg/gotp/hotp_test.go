package gotp

import (
	"testing"
)

var hotp = NewDefaultHOTP("4S62BZNFXXSZLCRO")

func TestHOTP_At(t *testing.T) {
	otp, err := hotp.At(12345)
	if err != nil {
		t.Fatalf("HOTP At error: %v", err)
	}
	if otp != "194001" {
		t.Errorf("HOTP generate otp error, got %s", otp)
	}
}

func TestHOTP_Verify(t *testing.T) {
	ok, err := hotp.Verify("194001", 12345)
	if err != nil {
		t.Fatalf("HOTP Verify error: %v", err)
	}
	if !ok {
		t.Error("verify failed")
	}
}
