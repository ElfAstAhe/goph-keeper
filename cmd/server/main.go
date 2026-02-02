package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ElfAstAhe/goph-keeper/internal/app/bootstrap"
	"github.com/ElfAstAhe/goph-keeper/internal/app/config"
)

func main() {
	printOrNA("Build version: %s\n", config.Version)
	printOrNA("Build date: %s\n", config.BuildTime)
	printOrNA("Project stage: %s\n", config.Stage)

	// app instance
	app := bootstrap.NewApp()
	//	defer app.Close()
	logger := app.GetLogger().GetLogger("main")
	//	defer _utl.CloseOnly(logger.(io.Closer))

	// app initialization
	logger.Info("app init")
	if err := app.Init(); err != nil {
		logger.Errorf("app initialization failed [%v]", err)
		defer app.Close()

		panic(errors.New("app initialization failed"))
	}

	// app run
	logger.Info("app run")
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Cancel()
		logger.Errorf("app run error [%v]", err)
	}

	//app.WG.Wait()

	// app close
	logger.Info("app close")
	if err := app.Close(); err != nil {
		logger.Errorf("app close error [%v]", err)

		panic(errors.New("app close failed"))
	}

	logger.Info("app shutdown")
}

func printOrNA(template, val string) {
	if strings.TrimSpace(val) == "" {
		fmt.Printf(template, "N/A")

		return
	}

	fmt.Printf(template, val)
}
