package config

import (
	"fmt"
	"strings"
)

// Флаги приложения
const (
	// FlagHTTPAddress - адрес и порт http listener, образец "localhost:8080" или ":8080"
	FlagHTTPAddress string = "http-address"

	// FlagHTTPSPrivateKeyPath - путь к файлу приватного ключа https
	FlagHTTPSPrivateKeyPath string = "https-private-key-path"

	// FlagHTTPSCertificatePath - путь к файлу сертификата https
	FlagHTTPSCertificatePath string = "https-certificate-path"

	// FlagHTTPSecure - используем https
	FlagHTTPSecure string = "http-secure"

	// FlagGRPCAddress - адрес и порт gRPC listener, образец "localhost:50051" или ":50051"
	FlagGRPCAddress string = "grpc-address"

	// FlagDatabaseDSN - строка соединения с БД
	FlagDatabaseDSN string = "database-dsn"

	// FlagDatabaseMaxOpenConnections - макс. кол-во соединений с БД
	FlagDatabaseMaxOpenConnections string = "database-max-open-connections"

	// FlagDatabaseMaxIdleConnections - макс. кол-во соединений с БД в простое
	FlagDatabaseMaxIdleConnections string = "database-max-idle-connections"

	// FlagDatabaseMaxIdleConnectionLifetime - макс. время жизни соединения с БД в секундах
	FlagDatabaseMaxIdleConnectionLifetime string = "database-max-idle-connection-lifetime"

	// FlagJWTSecretKey - серкетный ключ шифрования JWT
	FlagJWTSecretKey string = "jwt-secret-key"

	// FlagJWTLifetime - интервал протухания jwt в секундах
	FlagJWTLifetime string = "jwt-expiration-lifetime"

	// FlagCipherKey - ключ шифрования
	FlagCipherKey string = "cipher-key"
)

// Переменные среды
const (
	// EnvHTTPAddress - адрес и порт http listener, образец "localhost:8080" или ":8080"
	EnvHTTPAddress string = "GOPHKEEPER_HTTP_ADDRESS"

	// EnvHTTPSPrivateKeyPath - путь к файлу приватного ключа https
	EnvHTTPSPrivateKeyPath string = "GOPHKEEPER_HTTPS_PRIVATE_KEY_PATH"

	// EnvHTTPSCertificatePath - путь к файлу сертификата https
	EnvHTTPSCertificatePath string = "GOPHKEEPER_HTTPS_CERT_PATH"

	// EnvHTTPSecure - используем https
	EnvHTTPSecure string = "GOPHKEEPER_HTTP_SECURE"

	// EnvGRPCAddress - адрес и порт gRPC listener, образец "localhost:50051" или ":50051"
	EnvGRPCAddress string = "GOPHKEEPER_GRPC_ADDRESS"

	// EnvDatabaseDSN - строка соединения с БД
	EnvDatabaseDSN string = "GOPHKEEPER_DATABASE_DSN"

	// EnvDatabaseMaxOpenConnections - макс. кол-во соединений с БД
	EnvDatabaseMaxOpenConnections string = "GOPHKEEPER_DATABASE_MAX_OPEN_CONNECTIONS"

	// EnvDatabaseMaxIdleConnections - макс. время жизни соединения с БД в секундах
	EnvDatabaseMaxIdleConnections string = "GOPHKEEPER_DATABASE_MAX_IDLE_CONNECTIONS"

	// EnvDatabaseMaxIdleConnectionLifetime - макс. время жизни соединения с БД в секундах
	EnvDatabaseMaxIdleConnectionLifetime string = "GOPHKEEPER_DATABASE_MAX_IDLE_CONNECTION_LIFETIME"

	// EnvJWTSecretKey - серкетный ключ шифрования JWT
	EnvJWTSecretKey string = "GOPHKEEPER_JWT_SECRET_KEY"

	// EnvJWTLifetime - интервал протухания jwt в секундах
	EnvJWTLifetime string = "GOPHKEEPER_JWT_LIFETIME"

	// EnvCipherKey - ключ шифрования
	EnvCipherKey string = "GOPHKEEPER_CIPHER_KEY"
)

// Настройки по умолчанию
const (
	DefaultHTTPAddress                     string = ":8080"
	DefaultGRPCAddress                     string = ":9090"
	DefaultDatabaseMaxOpenConnections      int    = 20
	DefaultDatabaseMaxIdleConnections      int    = 5
	DefaultDatabaseIdleConnectionsLifeTime int    = 30
	DefaultJWTLifetime                     int    = 1800
	DefaultUseHTTPS                        bool   = false
)

const (
	undef string = "undefined"
)

// информация о версии приложения
var (
	// AppName - application name
	AppName string
	// Version - application version
	Version string
	// BuildTime - application build time
	BuildTime string
	// Stage - application stage (prod/dev/etc)
	Stage string
)

// BuildVersionInfo - build version info string
func BuildVersionInfo() string {
	if strings.HasPrefix(Version, "v") {
		Version = Version[1:]
	}

	return fmt.Sprintf("Application: %s, Version: v%s, BuildTime: %s, Stage: %s", AppName, Version, BuildTime, Stage)
}

func init() {
	if AppName == "" {
		AppName = "goph-keeper"
	}
	//if Version == "" {
	//	Version = undef
	//}
	//if BuildTime == "" {
	//	BuildTime = undef
	//}
	//if Stage == "" {
	//	Stage = undef
	//}
}
