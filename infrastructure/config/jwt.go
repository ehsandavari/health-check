package config

type SJwt struct {
	Algorithm       string `validate:"required"`
	PublicKey       string `validate:"required"`
	PrivateKey      string `validate:"required"`
	ExpiresAtMinute int    `validate:"required"`
}
