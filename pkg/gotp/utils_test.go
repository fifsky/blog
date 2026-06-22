package gotp

import (
	"testing"
)

func TestBuildUri(t *testing.T) {
	s, err := BuildUri(
		"totp",
		"4S62BZNFXXSZLCRO",
		"xlzd",
		"",
		"sha1",
		0,
		6,
		0,
	)
	if err != nil {
		t.Fatalf("BuildUri error: %v", err)
	}
	expected := "otpauth://totp/xlzd?secret=4S62BZNFXXSZLCRO"
	if s != expected {
		t.Errorf("BuildUri test failed.\n\texpected: %s,\n\tactual: %s", expected, s)
	}

	s2, err := BuildUri(
		"totp",
		"4S62BZNFXXSZLCRO",
		"xlzd",
		"Some Org",
		"sha1",
		0,
		6,
		0,
	)
	if err != nil {
		t.Fatalf("BuildUri error: %v", err)
	}
	expected2 := "otpauth://totp/Some%20Org:xlzd?issuer=Some%20Org&secret=4S62BZNFXXSZLCRO"
	if s2 != expected2 {
		t.Errorf("BuildUri test failed.\n\texpected: %s,\n\tactual: %s", expected2, s2)
	}

	s3, err := BuildUri(
		"hotp",
		"XXSZLCRO4S62BZNF",
		"mergenchik@gmail.com",
		"github.com",
		"sha1",
		0,
		6,
		0)
	if err != nil {
		t.Fatalf("BuildUri error: %v", err)
	}
	expected3 := "otpauth://hotp/github.com:mergenchik%40gmail.com?counter=0&issuer=github.com&secret=XXSZLCRO4S62BZNF"
	if s3 != expected3 {
		t.Errorf("BuildUri test failed.\n\texpected: %s,\n\tactual: %s", expected3, s3)
	}

	// Test invalid otp type
	_, err = BuildUri("invalid", "secret", "account", "", "sha1", 0, 6, 0)
	if err == nil {
		t.Error("BuildUri should return error for invalid otp type")
	}
}

func TestITob(t *testing.T) {
	var i int64 = 1524486261
	expect := []byte{0, 0, 0, 0, 90, 221, 208, 117}

	if string(expect) != string(itob(i)) {
		t.Error("itob error")
	}
}

func TestRandomSecret(t *testing.T) {
	secret, err := RandomSecret(64)
	if err != nil {
		t.Fatalf("RandomSecret error: %v", err)
	}
	if len(secret) == 0 {
		t.Error("RandomSecret returned empty secret")
	}
}

func TestIsSecretValid(t *testing.T) {
	valid, err := RandomSecret(64)
	if err != nil {
		t.Fatalf("RandomSecret error: %v", err)
	}
	if !IsSecretValid(valid) {
		t.Error("IsSecretValid error - RandomSecret(64) is not valid")
	}
	invalid := "asdsada"
	if IsSecretValid(invalid) {
		t.Error("IsSecretValid error - Bad secret is valid")
	}
}
