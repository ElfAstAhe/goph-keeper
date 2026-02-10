package main

import (
	"fmt"

	"github.com/ElfAstAhe/goph-keeper/internal/app/client/bootstrap"
	"github.com/ElfAstAhe/goph-keeper/internal/app/client/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func main() {
	fmt.Println(config.BuildVersionInfo())

	// app instance
	app := bootstrap.NewApp()

	// log
	log := app.GetLogger().GetLogger("main")

	// app init
	log.Info("app init")
	if err := app.Init(); err != nil {
		log.Errorf("app init failed [%v]", err)

		panic(errs.NewAppCommonError("app initialization failed", err))
	}

	// app run
	log.Info("app run")
	if err := app.Run(); err != nil {
		log.Errorf("app run error [%v]", err)

		panic(errs.NewAppCommonError("app run failed", err))
	}

	// app close
	log.Info("app close")
	if err := app.Close(); err != nil {
		log.Errorf("app close failed [%v]", err)
	}
}
