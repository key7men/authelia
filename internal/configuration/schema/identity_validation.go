package schema

import (
	"time"
)

// IdentityValidation represents the configuration for identity verification actions/flows.
type IdentityValidation struct {
	ResetPassword IdentityValidationResetPassword `koanf:"reset_password" json:"reset_password" jsonschema:"title=Reset Password" jsonschema_description:"Identity validation options for the Reset Password flow"`
	Credentials   IdentityValidationCredentials   `koanf:"credential_registration" json:"credential_registration" jsonschema:"title=Credential Registration" jsonschema_description:"Identity validation options for the Credential Modification flows"`
}

// IdentityValidationResetPassword represents the tunable aspects of the reset password identity verification action/flow.
type IdentityValidationResetPassword struct {
	JWTSecret       string        `koanf:"jwt_secret"  json:"jwt_secret" jsonschema:"title=JWT Secret" jsonschema_description:"The JWT secret used to sign the Reset Password flow JWT's"`
	EmailExpiration time.Duration `koanf:"email_expiration" json:"" jsonschema:"title=" jsonschema_description:""`
}

// IdentityValidationCredentials represents the tunable aspects of the credential control identity verification action/flow.
type IdentityValidationCredentials struct {
	EmailExpiration     time.Duration `koanf:"email_expiration" json:"email_expiration" jsonschema:"title=" jsonschema_description:"Duration of time the OTP code is considered valid"`
	ElevationExpiration time.Duration `koanf:"elevation_expiration" json:"elevation_expiration" jsonschema:"title=Elevation Expiration" jsonschema_description:"Duration of time the elevation can exist for after the user performs the validation"`
	OTPCharacters       int           `koanf:"otp_characters" json:"otp_characters" jsonschema:"title=OTP Characters" jsonschema_description:"Number of characters in the generated OTP codes"`
	Skip2FA             bool          `koanf:"skip_2fa" json:"skip_2fa" jsonschema:"title=Skip 2FA" jsonschema_description:"Skips 2FA requirement TODO"`
}
