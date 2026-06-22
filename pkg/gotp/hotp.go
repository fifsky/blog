package gotp

// HOTP implements HMAC-based one-time password counters (RFC 4226).
type HOTP struct {
	OTP
}

// NewHOTP creates a new HOTP instance with the given secret, digits, and hasher.
func NewHOTP(secret string, digits int, hasher *Hasher) *HOTP {
	otp := NewOTP(secret, digits, hasher)
	return &HOTP{OTP: otp}
}

// NewDefaultHOTP creates a HOTP instance with 6 digits and SHA1 hasher.
func NewDefaultHOTP(secret string) *HOTP {
	return NewHOTP(secret, 6, nil)
}

// At generates the OTP for the given counter value.
func (h *HOTP) At(count int) (string, error) {
	return h.generateOTP(int64(count))
}

// Verify checks whether the provided OTP matches the OTP generated at the given counter.
func (h *HOTP) Verify(otp string, count int) (bool, error) {
	generated, err := h.At(count)
	if err != nil {
		return false, err
	}
	return otp == generated, nil
}

// ProvisioningUri returns the provisioning URI for the OTP.
// This can be encoded in a QR Code and used to provision an OTP app like Google Authenticator.
//
// See also: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func (h *HOTP) ProvisioningUri(accountName, issuerName string, initialCount int) (string, error) {
	return BuildUri(
		OtpTypeHotp,
		h.secret,
		accountName,
		issuerName,
		h.hasher.HashName,
		initialCount,
		h.digits,
		0)
}
