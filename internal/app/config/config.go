package config

import "time"

const (
	DefaultJWTExpirationDuration time.Duration = 2 * time.Hour
)

// Config - конфигурация приложения
type Config struct {
	HTTPAddress           string `json:"http_address"`
	HTTPSPrivateKeyPath   string `json:"https_private_key_path,omitempty"`
	HTTPSCertificatePath  string `json:"https_certificate_path,omitempty"`
	GRPCAddress           string `json:"grpc_address"`
	DatabaseDSN           string `json:"database_dsn"`
	JWTSecretKey          string
	JWTExpirationDuration time.Duration `json:"jwt_expiration_duration"`
	CipherKey             string
}

func NewConfig() *Config {
	return &Config{}
}
