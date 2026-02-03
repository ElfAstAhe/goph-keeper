package config

import (
	"fmt"
	"strings"
)

// Флаги приложения
const (
	// FlagHTTPAddress - адрес и порт http listener, образец "localhost:8080" или ":8080"
	FlagHTTPAddress string = "a"
	// FlagHTTPSPrivateKeyPath - путь к файлу приватного ключа https
	FlagHTTPSPrivateKeyPath string = "https-private-key-path"
	// FlagHTTPSCertificatePath - путь к файлу сертификата https
	// FlagGRPCAddress - адрес и порт gRPC listener, образец "localhost:50051" или ":50051"
	FlagHTTPSCertificatePath string = "https-certificate-path"
	// FlagDatabaseDSN - строка соединения с БД
	FlagDatabaseDSN string = "d"
	// FlagJWTSecretKey - серкетный ключ шифрования JWT
	FlagJWTSecretKey string = "jwt-secret-key"
	// FlagCipherKey - ключ шифрования
	FlagCipherKey string = "cipher-key"
)

const (
	undef string = "undefined"
)

// информация о версии приложения
var (
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

	return fmt.Sprintf("Version: v%s, BuildTime: %s, Stage: %s", Version, BuildTime, Stage)
}

func init() {
	if Version == "" {
		Version = undef
	}
	if BuildTime == "" {
		BuildTime = undef
	}
	if Stage == "" {
		Stage = undef
	}
}
