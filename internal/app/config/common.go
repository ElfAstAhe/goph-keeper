package config

import (
	"fmt"
	"strings"
)

const undef string = "undefined"

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
