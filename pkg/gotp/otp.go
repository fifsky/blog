package gotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"hash"
	"math"
	"strings"
)

// Hasher represents the hash algorithm used for HMAC-based OTP generation.
type Hasher struct {
	HashName string
	Digest   func() hash.Hash
}

// OTP is the base struct for both HOTP and TOTP implementations.
// It holds the shared secret, digit count, and hash algorithm configuration.
type OTP struct {
	secret string  // secret in base32 format
	digits int     // number of integers in the OTP. Some apps expect this to be 6 digits, others support more.
	hasher *Hasher // digest function to use in the HMAC (expected to be sha1)
}

// NewOTP creates a new OTP instance with the given secret, digits, and hasher.
// If hasher is nil, SHA1 is used as the default.
func NewOTP(secret string, digits int, hasher *Hasher) OTP {
	if hasher == nil {
		hasher = &Hasher{
			HashName: "sha1",
			Digest:   sha1.New,
		}
	}
	return OTP{
		secret: secret,
		digits: digits,
		hasher: hasher,
	}
}

// generateOTP generates an OTP for the given input value.
// For HOTP, input is the counter value; for TOTP, input is the timecode.
func (o *OTP) generateOTP(input int64) (string, error) {
	if input < 0 {
		return "", fmt.Errorf("input must be positive integer")
	}
	secretBytes, err := o.byteSecret()
	if err != nil {
		return "", err
	}
	hasher := hmac.New(o.hasher.Digest, secretBytes)
	hasher.Write(itob(input))
	hmacHash := hasher.Sum(nil)

	offset := int(hmacHash[len(hmacHash)-1] & 0xf)
	code := ((int(hmacHash[offset]) & 0x7f) << 24) |
		((int(hmacHash[offset+1] & 0xff)) << 16) |
		((int(hmacHash[offset+2] & 0xff)) << 8) |
		(int(hmacHash[offset+3]) & 0xff)

	code = code % int(math.Pow10(o.digits))
	return fmt.Sprintf(fmt.Sprintf("%%0%dd", o.digits), code), nil
}

// byteSecret decodes the base32-encoded secret into raw bytes.
// It handles missing padding without modifying the original secret field.
func (o *OTP) byteSecret() ([]byte, error) {
	secret := o.secret
	missingPadding := len(secret) % 8
	if missingPadding != 0 {
		secret = secret + strings.Repeat("=", 8-missingPadding)
	}
	bytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("decode secret failed: %w", err)
	}
	return bytes, nil
}
