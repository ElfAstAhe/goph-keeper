package config

import (
	"fmt"
	"strings"
)

// общие опции
const (
	FlagOptConfig string = "config"
)

// команды первого запуска
const (
	FlagCmdGenConfig string = "gen-config"
	FlagCmdRegister  string = "register"
)

// команды профиля и аутентификации
const (
	FlagCmdProfile        string = "profile"
	FlagCmdGenKeys        string = "gen-keys"
	FlagCmdChangePassword string = "change-password"
)

// команды работы с файлами
const (
	FlagCmdGet    string = "get"
	FlagCmdPut    string = "put"
	FlagCmdRemove string = "remove"
	FlagCmdList   string = "list"
)

const (
	FlagOptAddress     string = "address"
	FlagOptUsername    string = "username"
	FlagOptPassword    string = "password"
	FLagOptOldPassword string = "old-password"
	FlagOptNewPassword string = "new-password"
	FlagOptDataKind    string = "data-kind"
	FlagOptName        string = "name"
	FlagOptData        string = "data"
	FlagOptPath        string = "path"
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
		AppName = "goph-keeper-client"
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
