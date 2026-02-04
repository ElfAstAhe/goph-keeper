package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElfAstAhe/goph-keeper/internal/app/server/bootstrap"
	"github.com/ElfAstAhe/goph-keeper/internal/app/server/config"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

func main() {
	fmt.Println(config.BuildVersionInfo())

	// app instance
	app := bootstrap.NewApp()

	// log
	log := app.GetLogger().GetLogger("main")

	// app initialization
	log.Info("app init")
	if err := app.Init(); err != nil {
		log.Errorf("app init failed [%v]", err)
		defer app.Close()

		panic(errs.NewAppCommonError("app initialization failed", err))
	}

	// app run
	log.Info("app run")
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Stop()
		log.Errorf("app run error [%v]", err)
	}

	app.WaitForStop()

	// app close
	log.Info("app close")
	if err := app.Close(); err != nil {
		log.Errorf("app close error [%v]", err)

		panic(errs.NewAppCommonError("app close failed", err))
	}

	log.Info("app shutdown")
}

func printOrNA(template, val string) {
	if strings.TrimSpace(val) == "" {
		fmt.Printf(template, "N/A")

		return
	}

	fmt.Printf(template, val)
}
