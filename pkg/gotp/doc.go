// Package gotp implements HMAC-based (HOTP) and Time-based (TOTP) one-time password
// algorithms as specified in RFC 4226 and RFC 6238.
//
// It can be used to implement two-factor (2FA) or multi-factor (MFA) authentication
// methods in applications that require users to log in.
//
// Basic usage:
//
//	totp := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO")
//	otp, err := totp.Now()
//
//	hotp := gotp.NewDefaultHOTP("4S62BZNFXXSZLCRO")
//	otp, err := hotp.At(0)
package gotp
