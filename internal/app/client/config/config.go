// Package config - application settings
// Приоритет загрузки конфига (от наивысшего к наименьшему)
//   - CLI
//   - file
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type AppCommands struct {
	GenConfig bool
	Register  bool

	Profile        bool
	GenKeys        bool
	ChangePassword bool

	Get    bool
	Put    bool
	Remove bool
	List   bool
}

func NewEmptyAppCommands() *AppCommands {
	return &AppCommands{}
}

func (c *AppCommands) Validate() error {
	cnt := c.cmdCount(c.GenConfig, c.Register, c.Profile, c.GenKeys, c.ChangePassword, c.Get, c.Put, c.Remove, c.List)
	// наличие хотя бы одной команды
	if cnt < 1 {
		return errs.NewAppConfigItemError("no commands", nil)
	}
	// только 1 команда
	if cnt > 1 {
		return errs.NewAppConfigItemError(fmt.Sprintf("to many commands: [%d]", cnt), nil)
	}

	return nil
}

func (c *AppCommands) cmdCount(commands ...bool) int {
	cnt := 0
	for _, cmd := range commands {
		if cmd {
			cnt++
		}
	}

	return cnt
}

type AppOptions struct {
	ConfigPath  string
	Address     string
	Username    string
	Password    string
	OldPassword string
	NewPassword string
	DataKind    string
	Name        string
	Data        string
	Path        string
}

func NewEmptyAppOptions() *AppOptions {
	return &AppOptions{}
}

type AppConfig struct {
	Address           string `json:"address"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
	PublicKey         string `json:"public_key"`
}

func NewEmptyAppConfig() *AppConfig {
	return &AppConfig{}
}

func (ac *AppConfig) Init(conf *AppConfig) {
	ac.Address = conf.Address
	ac.Username = conf.Username
	ac.EncryptedPassword = conf.EncryptedPassword
	ac.PublicKey = conf.PublicKey
}

type AppSettings struct {
	loadedConf *AppConfig
	conf       *AppConfig
	cmd        *AppCommands
	opts       *AppOptions
	fs         *flag.FlagSet
	keysHelper *utils.RSAKeysHelper
}

func NewAppSettings(keysHelper *utils.RSAKeysHelper) *AppSettings {
	return &AppSettings{
		loadedConf: NewEmptyAppConfig(),
		conf:       NewEmptyAppConfig(),
		cmd:        NewEmptyAppCommands(),
		opts:       NewEmptyAppOptions(),
		keysHelper: keysHelper,
	}
}

func (as *AppSettings) GetCmd() *AppCommands {
	return as.cmd
}

func (as *AppSettings) GetOpts() *AppOptions {
	return as.opts
}

func (as *AppSettings) GetConfig() *AppConfig {
	return as.conf
}

func (as *AppSettings) initCli() {
	as.fs = flag.NewFlagSet("config", flag.PanicOnError)

	as.fs.BoolVar(&as.cmd.GenConfig, FlagCmdGenConfig, false, "сгенерировать конфиг")
	as.fs.BoolVar(&as.cmd.Register, FlagCmdRegister, false, "регистрация в сервисе")

	as.fs.BoolVar(&as.cmd.Profile, FlagCmdProfile, false, "информация о профиле")
	as.fs.BoolVar(&as.cmd.GenKeys, FlagCmdGenKeys, false, "сгенерировать новую пару RSA ключей")
	as.fs.BoolVar(&as.cmd.ChangePassword, FlagCmdChangePassword, false, "сменить пароль")

	as.fs.BoolVar(&as.cmd.Get, FlagCmdGet, false, "получить данные")
	as.fs.BoolVar(&as.cmd.Put, FlagCmdPut, false, "сохранить данные")
	as.fs.BoolVar(&as.cmd.Remove, FlagCmdRemove, false, "удалить данные")
	as.fs.BoolVar(&as.cmd.List, FlagCmdList, false, "список данных")

	as.fs.StringVar(&as.opts.ConfigPath, FlagOptConfig, "", "файл конфига")
	as.fs.StringVar(&as.opts.Address, FlagOptAddress, "", "хост:порт (http://example.org:8080)")
	as.fs.StringVar(&as.opts.Username, FlagOptUsername, "", "имя пользователя")
	as.fs.StringVar(&as.opts.Password, FlagOptPassword, "", "пароль")
	as.fs.StringVar(&as.opts.OldPassword, FLagOptOldPassword, "", "старый пароль, обязан совпадать с текущим паролем")
	as.fs.StringVar(&as.opts.NewPassword, FlagOptNewPassword, "", "новый пароль, не пустой, не совпадает со старым")
	as.fs.StringVar(&as.opts.DataKind, FlagOptDataKind, "", "тип сохранённых данных (credential,plaintext,binary,bankcard)")
	as.fs.StringVar(&as.opts.Name, FlagOptName, "", "намиенование данных")
	as.fs.StringVar(&as.opts.Data, FlagOptData, "", "текстовые данные")
	as.fs.StringVar(&as.opts.Path, FlagOptPath, "", "путь к файлу конфига")
}

func (as *AppSettings) loadCli() (err error) {
	// init
	as.initCli()

	// act
	if err = as.fs.Parse(os.Args[1:]); err != nil {
		return errs.NewAppConfigError("parse cli flags", err)
	}
	defer func() {
		if r := recover(); r != nil {
			// Проверяем, является ли r ошибкой
			recoveryErr, ok := r.(error)
			if !ok {
				// Если это строка или что-то другое, приводим к виду error вручную
				recoveryErr = fmt.Errorf("%v", r)
			}
			err = errs.NewAppConfigError("parse cli flags panic", recoveryErr)
		}
	}()

	// result
	return nil
}

func (as *AppSettings) loadFile() error {
	var err error
	if as.opts.ConfigPath == "" {
		as.opts.ConfigPath, err = as.buildDefaultConfigPath()
		if err != nil {
			return errs.NewAppConfigError("build default config path", err)
		}
	}
	f, err := os.OpenFile(as.opts.ConfigPath, os.O_RDONLY, 0600)
	if err != nil {
		return errs.NewAppConfigError("open config file", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(as.loadedConf); err != nil {
		return errs.NewAppConfigError("decode config JSON", err)
	}

	return nil
}

func (as *AppSettings) buildDefaultConfigPath() (string, error) {
	// 1. Получаем путь к запущенному бинарнику (например, /home/user/goph-keeper/bin/server)
	exePath, err := os.Executable()
	if err != nil {
		return "", errs.NewAppConfigError("get executable path", err)
	}

	// 2. Берем директорию бинарника (/home/user/goph-keeper/bin)
	exeDir := filepath.Dir(exePath)

	// 4. Собираем путь (например, /home/user/goph-keeper/config/config.yaml)
	configPath := filepath.Join(exeDir, "goph-keeper-client-conf.json")

	return configPath, nil
}

func (as *AppSettings) mergeConfig() error {
	var err error

	// инициализируем
	as.conf.Init(as.loadedConf)

	// данные
	if as.opts.Address != "" {
		as.conf.Address = as.opts.Address
	}
	if as.opts.Username != "" {
		as.conf.Username = as.opts.Username
	}
	if as.opts.Password != "" {
		as.conf.EncryptedPassword, err = as.encryptPassword(as.opts.Password, as.conf.PublicKey)
		if err != nil {
			return errs.NewAppConfigError("encrypt password", err)
		}
	}

	return nil
}

func (as *AppSettings) encryptPassword(password string, publicKey string) (string, error) {
	pubKey, err := as.keysHelper.ParsePublicKey(publicKey)
	if err != nil {
		return "", errs.NewAppConfigError("parse public key", err)
	}

	encryptedPassword, err := as.keysHelper.EncryptString(password, pubKey)
	if err != nil {
		return "", errs.NewAppConfigError("encrypt password", err)
	}

	return encryptedPassword, nil
}

func (as *AppSettings) Load() error {
	err := as.loadCli()
	if err != nil {
		return errs.NewAppConfigError("load cli", err)
	}

	_ = as.loadFile()

	return as.mergeConfig()
}

func (as *AppSettings) SaveFile() error {
	f, err := os.OpenFile(as.opts.ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errs.NewAppCommonError(fmt.Sprintf("open file at [%s]", as.opts.ConfigPath), err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(as.conf)
}

func (as *AppSettings) Validate() error {
	if err := as.cmd.Validate(); err != nil {
		return errs.NewAppConfigError("validate cmd", err)
	}

	if as.conf.Address == "" {
		return errs.NewAppConfigItemError("address empty", nil)
	}
	if as.conf.Username == "" {
		return errs.NewAppConfigItemError("username empty", nil)
	}
	if as.conf.PublicKey == "" {
		return errs.NewAppConfigItemError("public key empty", nil)
	}

	return nil
}
