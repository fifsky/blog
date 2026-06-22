package gotp

import (
	"testing"
	"time"
)

var totp = NewDefaultTOTP("4S62BZNFXXSZLCRO")

func TestTOTP_At(t *testing.T) {
	now, err := totp.Now()
	if err != nil {
		t.Fatalf("TOTP Now error: %v", err)
	}
	at, err := totp.At(currentTimestamp())
	if err != nil {
		t.Fatalf("TOTP At error: %v", err)
	}
	if now != at {
		t.Error("TOTP generate otp error!")
	}
}

func TestTOTP_AtTime(t *testing.T) {
	now, err := totp.Now()
	if err != nil {
		t.Fatalf("TOTP Now error: %v", err)
	}
	atTime, err := totp.AtTime(time.Now())
	if err != nil {
		t.Fatalf("TOTP AtTime error: %v", err)
	}
	if now != atTime {
		t.Error("TOTP at time generate otp error!")
	}
}

func TestTOTP_NowWithExpiration(t *testing.T) {
	otp, exp, err := totp.NowWithExpiration()
	if err != nil {
		t.Fatalf("TOTP NowWithExpiration error: %v", err)
	}
	cts := currentTimestamp()
	now, err := totp.Now()
	if err != nil {
		t.Fatalf("TOTP Now error: %v", err)
	}
	if otp != now {
		t.Error("TOTP generate otp error!")
	}
	atCts, err := totp.At(cts + 30)
	if err != nil {
		t.Fatalf("TOTP At error: %v", err)
	}
	atExp, err := totp.At(exp)
	if err != nil {
		t.Fatalf("TOTP At error: %v", err)
	}
	if atCts != atExp {
		t.Error("TOTP expiration otp error!")
	}
}

func TestTOTP_Verify(t *testing.T) {
	ok, err := totp.Verify("179394", 1524485781)
	if err != nil {
		t.Fatalf("TOTP Verify error: %v", err)
	}
	if !ok {
		t.Error("verify failed")
	}
}

func TestTOTP_VerifyWithWindow(t *testing.T) {
	// Generate OTP at timestamp 1524485781
	otp, err := totp.At(1524485781)
	if err != nil {
		t.Fatalf("TOTP At error: %v", err)
	}

	// Verify with window=1 should match at timestamp ±30 seconds
	ok, err := totp.VerifyWithWindow(otp, 1524485781+30, 1)
	if err != nil {
		t.Fatalf("TOTP VerifyWithWindow error: %v", err)
	}
	if !ok {
		t.Error("VerifyWithWindow should match with window=1")
	}

	ok, err = totp.VerifyWithWindow(otp, 1524485781-30, 1)
	if err != nil {
		t.Fatalf("TOTP VerifyWithWindow error: %v", err)
	}
	if !ok {
		t.Error("VerifyWithWindow should match with window=1")
	}

	// Verify with window=0 should NOT match at timestamp ±30 seconds
	ok, err = totp.VerifyWithWindow(otp, 1524485781+30, 0)
	if err != nil {
		t.Fatalf("TOTP VerifyWithWindow error: %v", err)
	}
	if ok {
		t.Error("VerifyWithWindow should not match with window=0")
	}
}

func TestTOTP_ProvisioningUri(t *testing.T) {
	expect := "otpauth://totp/github:xlzd?issuer=github&secret=4S62BZNFXXSZLCRO"
	uri, err := totp.ProvisioningUri("xlzd", "github")
	if err != nil {
		t.Fatalf("ProvisioningUri error: %v", err)
	}
	if expect != uri {
		t.Errorf("ProvisioningUri error.\n\texpected: %s,\n\tactual: %s", expect, uri)
	}
}
