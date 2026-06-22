# GOTP - The Golang One-Time Password Library

GOTP is a Golang package for generating and verifying one-time passwords. It can be used to implement two-factor (2FA) or multi-factor (MFA) authentication methods in anywhere that requires users to log in.

Open MFA standards are defined in [RFC 4226][RFC 4226] (HOTP: An HMAC-Based One-Time Password Algorithm) and in [RFC 6238][RFC 6238] (TOTP: Time-Based One-Time Password Algorithm). GOTP implements server-side support for both of these standards.

### Time-based OTPs

```go
totp := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO")
otp, err := totp.Now()           // current otp, e.g. "123456"
otp, err := totp.At(1524486261)  // otp of timestamp 1524486261, e.g. "123456"

// OTP verified for a given timestamp
ok, err := totp.Verify("492039", 1524486261)  // true
ok, err := totp.Verify("492039", 1520000000)  // false

// OTP verified with time window tolerance (handles clock drift)
ok, err := totp.VerifyWithWindow("492039", 1524486261, 1)  // true, allows ±1 interval

// generate a provisioning uri
totp.ProvisioningUri("demoAccountName", "issuerName")
// otpauth://totp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&issuer=issuerName
```

### Counter-based OTPs

```go
hotp := gotp.NewDefaultHOTP("4S62BZNFXXSZLCRO")
otp, err := hotp.At(0)  // e.g. "944181"
otp, err := hotp.At(1)  // e.g. "770975"

// OTP verified for a given counter
ok, err := hotp.Verify("944181", 0)  // true
ok, err := hotp.Verify("944181", 1)  // false

// generate a provisioning uri
hotp.ProvisioningUri("demoAccountName", "issuerName", 1)
// otpauth://hotp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&counter=1&issuer=issuerName
```

### Generate random secret

```go
secret, err := gotp.RandomSecret(16)  // e.g. "LMT4URYNZKEWZRAA"
```

### Google Authenticator Compatible

GOTP works with the Google Authenticator iPhone and Android app, as well as other OTP apps like Authy.
GOTP includes the ability to generate provisioning URIs for use with the QR Code
scanner built into these MFA client apps via `ProvisioningUri` method:

```go
gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").ProvisioningUri("demoAccountName", "issuerName")
// otpauth://totp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&issuer=issuerName

gotp.NewDefaultHOTP("4S62BZNFXXSZLCRO").ProvisioningUri("demoAccountName", "issuerName", 1)
// otpauth://hotp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&counter=1&issuer=issuerName
```

This URL can then be rendered as a QR Code which can then be scanned and added to the users list of OTP credentials.
