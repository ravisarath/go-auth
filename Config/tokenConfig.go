package Config

import "os"

type TokenConfig struct {
	Path     string
	MaxAge   int
	Domain   string
	Secure   bool
	HttpOnly bool
}

func BuildTokenConfig() *TokenConfig {
	if len(os.Getenv("path")) == 0 {
		os.Setenv("path", "/")
	}
	if len(os.Getenv("domain")) == 0 {
		os.Setenv("domain", "0.0.0.0")
	}
	tokenconf := TokenConfig{
		Path:     os.Getenv("path"),
		MaxAge:   3600,
		Domain:   os.Getenv("domain"),
		Secure:   true,
		HttpOnly: true,
	}
	return &tokenconf
}

type KeyConfig struct {
	AtSecret string
	RfSecret string
}

func BuildKeyConfig() *KeyConfig {
	if len(os.Getenv("ACCESS_SECRET")) == 0 {
		os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	}
	if len(os.Getenv("REFRESH_SECRET")) == 0 {
		os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	}
	keyconf := KeyConfig{
		AtSecret: os.Getenv("ACCESS_SECRET"),

		RfSecret: os.Getenv("REFRESH_SECRET"),
	}
	return &keyconf
}
