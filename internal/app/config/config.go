package config

import (
	"flag"
	"os"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/caarlos0/env/v6"
)

type HTTPConfig struct {
	Address        string `json:"address" env:"GOPHKEEPER_HTTP_ADDRESS"`
	PrivateKeyPath string `json:"private_key_path,omitempty" env:"GOPHKEEPER_HTTPS_PRIVATE_KEY_PATH"`
	CertPath       string `json:"cert_path,omitempty" env:"GOPHKEEPER_HTTPS_CERT_PATH"`
	UseHTTPS       bool   `json:"use_https,omitempty" env:"GOPHKEEPER_HTTP_SECURE"`
}

func NewDefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Address:        DefaultHTTPAddress,
		PrivateKeyPath: "",
		CertPath:       "",
		UseHTTPS:       DefaultUseHTTPS,
	}
}

func (ch *HTTPConfig) Validate() error {
	if ch.Address == "" {
		return errs.NewAppConfigItemError("http address", nil)
	}

	if ch.UseHTTPS {
		if ch.PrivateKeyPath == "" {
			return errs.NewAppConfigItemError("https private key path empty", nil)
		}
		if ch.CertPath == "" {
			return errs.NewAppConfigItemError("https certificate path empty", nil)
		}
	}

	return nil
}

func (ch *HTTPConfig) IsHTTPS() bool {
	return ch.UseHTTPS
}

type GRPCConfig struct {
	Address string `json:"address" env:"GOPHKEEPER_GRPC_ADDRESS"`
}

func NewDefaultGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		Address: DefaultGRPCAddress,
	}
}

func (gc *GRPCConfig) Validate() error {
	if gc.Address == "" {
		return errs.NewAppConfigItemError("gRPC address", nil)
	}

	return nil
}

type DatabaseConfig struct {
	DSN                       string `json:"dsn" env:"GOPHKEEPER_DATABASE_DSN"`
	MaxOpenConnections        int    `json:"max_open_connections,omitempty" env:"GOPHKEEPER_DATABASE_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections        int    `json:"max_idle_connections,omitempty" env:"GOPHKEEPER_DATABASE_MAX_IDLE_CONNECTIONS"`
	MaxIdleConnectionLifetime int    `json:"max_idle_connection_lifetime,omitempty" env:"GOPHKEEPER_DATABASE_MAX_IDLE_CONNECTION_LIFETIME"`
}

func NewDefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN:                       "",
		MaxOpenConnections:        DefaultDatabaseMaxOpenConnections,
		MaxIdleConnections:        DefaultDatabaseMaxIdleConnections,
		MaxIdleConnectionLifetime: DefaultDatabaseIdleConnectionsLifeTime,
	}
}

func (dc *DatabaseConfig) Validate() error {
	if dc.DSN == "" {
		return errs.NewAppConfigItemError("database DSN", nil)
	}
	if dc.MaxOpenConnections <= 0 {
		return errs.NewAppConfigItemError("database max open connections", nil)
	}
	if dc.MaxIdleConnections <= 0 {
		return errs.NewAppConfigItemError("database max idle connections", nil)
	}
	if dc.MaxIdleConnectionLifetime <= 0 {
		return errs.NewAppConfigItemError("database max idle connection lifetime", nil)
	}

	return nil
}

type JWTConfig struct {
	SecretKey string `env:"GOPHKEEPER_JWT_SECRET_KEY"`
	Lifetime  int    `json:"lifetime" env:"GOPHKEEPER_JWT_LIFETIME"`
}

func NewDefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: "",
		Lifetime:  DefaultJWTLifetime,
	}
}

func (ch *JWTConfig) Validate() error {
	if ch.SecretKey == "" {
		return errs.NewAppConfigItemError("jwt secret key", nil)
	}
	if ch.Lifetime <= 0 {
		return errs.NewAppConfigItemError("jwt lifetime seconds", nil)
	}

	return nil
}

// Config - конфигурация приложения
//
// load priority (from high to low):
//   - ENV Variables
//   - Flags
//   - Default values
type Config struct {
	HTTPConfig     *HTTPConfig     `json:"http_config"`
	GRPCConfig     *GRPCConfig     `json:"grpc_config"`
	DatabaseConfig *DatabaseConfig `json:"database_config"`
	JWTConfig      *JWTConfig      `json:"jwt_config"`
	CipherKey      string          `env:"GOPHKEEPER_CIPHER_KEY"`
}

func NewConfig() *Config {
	return &Config{
		HTTPConfig:     NewDefaultHTTPConfig(),
		GRPCConfig:     NewDefaultGRPCConfig(),
		DatabaseConfig: NewDefaultDatabaseConfig(),
		JWTConfig:      NewDefaultJWTConfig(),
		CipherKey:      "",
	}
}

// Load - загрузка конфигурации ведётся в обратном порядке:
//   - defaults
//   - cli
//   - env
func (c *Config) Load() error {
	// cli
	if err := c.loadCli(); err != nil {
		return errs.NewAppConfigError("load cli ", err)
	}

	// env
	if err := c.loadEnv(); err != nil {
		return errs.NewAppConfigError("load env", err)
	}

	return nil
}

// Validate - проверка конфигурационных значений
func (c *Config) Validate() error {
	// http
	if err := c.HTTPConfig.Validate(); err != nil {
		return errs.NewAppConfigError("validate http config", err)
	}
	// gRPC
	if err := c.GRPCConfig.Validate(); err != nil {
		return errs.NewAppConfigError("validate gRPC config", err)
	}
	// Database
	if err := c.DatabaseConfig.Validate(); err != nil {
		return errs.NewAppConfigError("validate database config", err)
	}
	// JWT
	if err := c.JWTConfig.Validate(); err != nil {
		return errs.NewAppConfigError("validate jwt config", err)
	}
	// cipher
	if c.CipherKey == "" {
		return errs.NewAppConfigError("validate cipher key empty", nil)
	}

	return nil
}

func (c *Config) initCli() *flag.FlagSet {
	cli := flag.NewFlagSet("config", flag.PanicOnError)

	cli.StringVar(&c.HTTPConfig.Address, FlagHTTPAddress, DefaultHTTPAddress, "http address")
	cli.StringVar(&c.HTTPConfig.PrivateKeyPath, FlagHTTPSPrivateKeyPath, "", "https private key path")
	cli.StringVar(&c.HTTPConfig.CertPath, FlagHTTPSCertificatePath, "", "https certificate path")
	cli.BoolVar(&c.HTTPConfig.UseHTTPS, FlagHTTPSecure, false, "use https instead of http")

	cli.StringVar(&c.GRPCConfig.Address, FlagGRPCAddress, DefaultGRPCAddress, "gRPC address")

	cli.StringVar(&c.DatabaseConfig.DSN, FlagDatabaseDSN, "", "database DSN")
	cli.IntVar(&c.DatabaseConfig.MaxOpenConnections, FlagDatabaseMaxOpenConnections, DefaultDatabaseMaxOpenConnections, "database max open connections")
	cli.IntVar(&c.DatabaseConfig.MaxIdleConnections, FlagDatabaseMaxIdleConnections, DefaultDatabaseMaxIdleConnections, "database max idle connections")
	cli.IntVar(&c.DatabaseConfig.MaxIdleConnectionLifetime, FlagDatabaseMaxIdleConnectionLifetime, DefaultDatabaseIdleConnectionsLifeTime, "database max idle connection lifetime seconds")

	cli.StringVar(&c.JWTConfig.SecretKey, FlagJWTSecretKey, "", "JWT secret key")
	cli.IntVar(&c.JWTConfig.Lifetime, FlagJWTLifetime, DefaultJWTLifetime, "JWT token lifetime seconds")

	cli.StringVar(&c.CipherKey, FlagCipherKey, "", "cipher key")

	return cli
}

func (c *Config) loadCli() error {
	// init
	flagSet := c.initCli()

	// act
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return errs.NewAppConfigError("load cli flags", err)
	}

	// special cases
	if c.cliFlagExists(FlagHTTPSecure, flagSet) {
		c.HTTPConfig.UseHTTPS = true
	}

	// result
	return nil
}

func (c *Config) cliFlagExists(name string, flagSet *flag.FlagSet) bool {
	res := false

	flagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			res = true
		}
	})

	return res
}

func (c *Config) loadEnv() error {
	// act
	err := env.Parse(c)
	if err != nil {
		return errs.NewAppConfigError("load env variables", err)
	}

	// special cases
	if c.envVarExists(EnvHTTPSecure) {
		c.HTTPConfig.UseHTTPS = true
	}

	// result
	return nil
}

func (c *Config) envVarExists(name string) bool {
	_, ok := os.LookupEnv(name)

	return ok
}
