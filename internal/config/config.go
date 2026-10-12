package config

import (
	"crypto/tls"
	"errors"
	"os"
)

type Config struct {
	APIPort       string
	Environment   string
	EncryptionKey []byte
	TLSCertPath   string
	TLSKeyPath    string
}

func Load() (*Config, error) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	secretKey := os.Getenv("APP_ENCRYPTION_KEY")
	if env == "production" && len(secretKey) != 32 {
		return nil, errors.New("secure encryption key must be exactly 32 bytes in production")
	}

	return &Config{
		APIPort:       port,
		Environment:   env,
		EncryptionKey: []byte(secretKey),
		TLSCertPath:   os.Getenv("TLS_CERT_PATH"),
		TLSKeyPath:    os.Getenv("TLS_KEY_PATH"),
	}, nil
}

func GetSecureTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}
}
