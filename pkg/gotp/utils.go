package gotp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	OtpTypeTotp = "totp"
	OtpTypeHotp = "hotp"
)

// otpPathEscape percent-encodes a string for use in the otpauth URI path.
// It extends url.PathEscape to also encode @ which is valid in generic URI paths
// but must be percent-encoded per the otpauth URI specification.
func otpPathEscape(s string) string {
	return strings.ReplaceAll(url.PathEscape(s), "@", "%40")
}

// BuildUri constructs the provisioning URI for the OTP; works for either TOTP or HOTP.
// This can then be encoded in a QR Code and used to provision the Google Authenticator app.
// For module-internal use.
//
// See also: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func BuildUri(otpType, secret, accountName, issuerName, algorithm string, initialCount, digits int, period int) (string, error) {
	if otpType != OtpTypeHotp && otpType != OtpTypeTotp {
		return "", fmt.Errorf("otp type error, got %s", otpType)
	}

	// Build label: "issuer:accountName" or just "accountName"
	// Both issuer and accountName need percent-encoding per otpauth URI spec.
	// url.PathEscape does not encode @ (valid in URI paths), but the otpauth
	// spec requires it to be percent-encoded, so we handle it manually.
	label := otpPathEscape(accountName)
	if issuerName != "" {
		label = otpPathEscape(issuerName) + ":" + label
	}

	// Build query parameters
	q := url.Values{}
	q.Set("secret", secret)
	if issuerName != "" {
		q.Set("issuer", issuerName)
	}
	if algorithm != "" && algorithm != "sha1" {
		q.Set("algorithm", strings.ToUpper(algorithm))
	}
	if digits != 0 && digits != 6 {
		q.Set("digits", fmt.Sprintf("%d", digits))
	}
	if period != 0 && period != 30 {
		q.Set("period", fmt.Sprintf("%d", period))
	}
	if otpType == OtpTypeHotp {
		q.Set("counter", fmt.Sprintf("%d", initialCount))
	}

	// Assemble URI string directly.
	// url.Values.Encode() uses + for spaces, but otpauth spec requires %20.
	rawQuery := strings.ReplaceAll(q.Encode(), "+", "%20")
	return fmt.Sprintf("otpauth://%s/%s?%s", otpType, label, rawQuery), nil
}

// currentTimestamp returns the current Unix timestamp.
func currentTimestamp() int64 {
	return time.Now().Unix()
}

// itob converts an integer to a byte array (big-endian, 8 bytes).
func itob(integer int64) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}

// RandomSecret generates a random secret of the given length (number of bytes).
// Returns a base32-encoded string without padding.
func RandomSecret(length int) (string, error) {
	secret := make([]byte, length)
	gen, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	if gen != length {
		return "", fmt.Errorf("failed to generate random secret: expected %d bytes, got %d", length, gen)
	}
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(secret), nil
}

// IsSecretValid checks whether a given base32 secret string is valid.
func IsSecretValid(secret string) bool {
	missingPadding := len(secret) % 8
	if missingPadding != 0 {
		secret = secret + strings.Repeat("=", 8-missingPadding)
	}
	_, err := base32.StdEncoding.DecodeString(secret)
	return err == nil
}
