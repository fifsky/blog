package gotp

import "time"

// TOTP implements time-based one-time password counters (RFC 6238).
type TOTP struct {
	OTP
	interval int
}

// NewTOTP creates a new TOTP instance with the given secret, digits, interval, and hasher.
func NewTOTP(secret string, digits, interval int, hasher *Hasher) *TOTP {
	otp := NewOTP(secret, digits, hasher)
	return &TOTP{OTP: otp, interval: interval}
}

// NewDefaultTOTP creates a TOTP instance with 6 digits, 30-second interval, and SHA1 hasher.
func NewDefaultTOTP(secret string) *TOTP {
	return NewTOTP(secret, 6, 30, nil)
}

// At generates the OTP for the given timestamp.
func (t *TOTP) At(timestamp int64) (string, error) {
	return t.generateOTP(t.timecode(timestamp))
}

// AtTime generates the OTP for the given time.Time.
func (t *TOTP) AtTime(timestamp time.Time) (string, error) {
	return t.At(timestamp.Unix())
}

// Now generates the current time OTP.
func (t *TOTP) Now() (string, error) {
	return t.At(currentTimestamp())
}

// NowWithExpiration generates the current time OTP along with its expiration timestamp.
func (t *TOTP) NowWithExpiration() (string, int64, error) {
	ts := currentTimestamp()
	timeCodeInt64 := t.timecode(ts)
	otp, err := t.generateOTP(timeCodeInt64)
	if err != nil {
		return "", 0, err
	}
	expirationTime := (timeCodeInt64 + 1) * int64(t.interval)
	return otp, expirationTime, nil
}

// Verify checks whether the provided OTP matches the OTP generated at the given timestamp.
func (t *TOTP) Verify(otp string, timestamp int64) (bool, error) {
	generated, err := t.At(timestamp)
	if err != nil {
		return false, err
	}
	return otp == generated, nil
}

// VerifyTime checks whether the provided OTP matches the OTP generated at the given time.Time.
func (t *TOTP) VerifyTime(otp string, timestamp time.Time) (bool, error) {
	return t.Verify(otp, timestamp.Unix())
}

// VerifyWithWindow checks whether the provided OTP matches the OTP generated at the given timestamp,
// allowing a window of validWindow time intervals before and after the timestamp.
// This is useful for handling clock drift between client and server.
func (t *TOTP) VerifyWithWindow(otp string, timestamp int64, validWindow int) (bool, error) {
	for i := -validWindow; i <= validWindow; i++ {
		result, err := t.Verify(otp, timestamp+int64(i*t.interval))
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

// ProvisioningUri returns the provisioning URI for the OTP.
// This can be encoded in a QR Code and used to provision an OTP app like Google Authenticator.
//
// See also: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func (t *TOTP) ProvisioningUri(accountName, issuerName string) (string, error) {
	return BuildUri(
		OtpTypeTotp,
		t.secret,
		accountName,
		issuerName,
		t.hasher.HashName,
		0,
		t.digits,
		t.interval)
}

func (t *TOTP) timecode(timestamp int64) int64 {
	return timestamp / int64(t.interval)
}
