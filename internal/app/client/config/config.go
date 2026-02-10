// Package config - application settings
// Приоритет загрузки конфига
//   - CLI
//   - file
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

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
	cmd               *AppCommands
	opts              *AppOptions
	fs                *flag.FlagSet
	keysHelper        *utils.RSAKeysHelper
}

func NewAppConfig(keysHelper *utils.RSAKeysHelper) *AppConfig {
	return &AppConfig{
		cmd:        NewEmptyAppCommands(),
		opts:       NewEmptyAppOptions(),
		keysHelper: keysHelper,
	}
}

func (ac *AppConfig) initCli() {
	ac.fs = flag.NewFlagSet("config", flag.PanicOnError)

	ac.fs.BoolVar(&ac.cmd.GenConfig, FlagCmdGenConfig, false, "сгенерировать конфиг")
	ac.fs.BoolVar(&ac.cmd.Register, FlagCmdRegister, false, "регистрация в сервисе")

	ac.fs.BoolVar(&ac.cmd.Profile, FlagCmdProfile, false, "информация о профиле")
	ac.fs.BoolVar(&ac.cmd.GenKeys, FlagCmdGenKeys, false, "сгенерировать новую пару RSA ключей")
	ac.fs.BoolVar(&ac.cmd.ChangePassword, FlagCmdChangePassword, false, "сменить пароль")

	ac.fs.BoolVar(&ac.cmd.Get, FlagCmdGet, false, "получить данные")
	ac.fs.BoolVar(&ac.cmd.Put, FlagCmdPut, false, "сохранить данные")
	ac.fs.BoolVar(&ac.cmd.Remove, FlagCmdRemove, false, "удалить данные")
	ac.fs.BoolVar(&ac.cmd.List, FlagCmdList, false, "список данных")

	ac.fs.StringVar(&ac.opts.ConfigPath, FlagOptConfig, "", "файл конфига")
	ac.fs.StringVar(&ac.opts.Address, FlagOptAddress, "", "хост:порт (http://example.org:8080)")
	ac.fs.StringVar(&ac.opts.Username, FlagOptUsername, "", "имя пользователя")
	ac.fs.StringVar(&ac.opts.Password, FlagOptPassword, "", "пароль")
	ac.fs.StringVar(&ac.opts.OldPassword, FLagOptOldPassword, "", "старый пароль, обязан совпадать с текущим паролем")
	ac.fs.StringVar(&ac.opts.NewPassword, FlagOptNewPassword, "", "новый пароль, не пустой, не совпадает со старым")
	ac.fs.StringVar(&ac.opts.DataKind, FlagOptDataKind, "", "тип сохранённых данных (credential,plaintext,binary,bankcard)")
	ac.fs.StringVar(&ac.opts.Name, FlagOptName, "", "намиенование данных")
	ac.fs.StringVar(&ac.opts.Data, FlagOptData, "", "текстовые данные")
	ac.fs.StringVar(&ac.opts.Path, FlagOptPath, "", "путь к файлу конфига")
}

func (ac *AppConfig) loadCli() (err error) {
	// init
	ac.initCli()

	// act
	if err = ac.fs.Parse(os.Args[1:]); err != nil {
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

func (ac *AppConfig) loadFile() error {
	if ac.opts.ConfigPath != "" {
		ac.opts.ConfigPath = ac.buildDefaultConfigPath()
	}
	f, err := os.OpenFile(ac.opts.ConfigPath, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(ac); err != nil {
		return err
	}

	return nil
}

func (ac *AppConfig) buildDefaultConfigPath() string {

}

func (ac *AppConfig) UpdateConfig(address string, username string, password string, publicKey string) error {
	var err error

	ac.Address = address
	ac.Username = username
	ac.PublicKey = publicKey
	ac.EncryptedPassword, err = ac.encryptPassword(password, ac.PublicKey)
	if err != nil {
		return errs.NewAppConfigError("encrypt password", err)
	}

	return nil
}

func (ac *AppConfig) encryptPassword(password string, publicKey string) (string, error) {
	pubKey, err := ac.keysHelper.ParsePublicKey(publicKey)
	if err != nil {
		return "", errs.NewAppConfigError("parse public key", err)
	}

	encryptedPassword, err := ac.keysHelper.EncryptString(password, pubKey)
	if err != nil {
		return "", errs.NewAppConfigError("encrypt password", err)
	}

	return encryptedPassword, nil
}

func (ac *AppConfig) Load() error {
	err := ac.loadCli()
	if err != nil {
		return errs.NewAppConfigError("load cli", err)
	}

	err = ac.loadFile()
	if err != nil {
		return errs.NewAppConfigError("load file", err)
	}

	return nil
}

func (ac *AppConfig) SaveFile() error {
	f, err := os.OpenFile(ac.opts.ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ac)
}

func (ac *AppConfig) Validate() error {
	if err := ac.cmd.Validate(); err != nil {
		return errs.NewAppConfigError("validate cmd", err)
	}

	if ac.Address == "" {
		return errs.NewAppConfigItemError("address empty", nil)
	}
	if ac.Username == "" {
		return errs.NewAppConfigItemError("username empty", nil)
	}
	if ac.PublicKey == "" {
		return errs.NewAppConfigItemError("public key empty", nil)
	}

	return nil
}
