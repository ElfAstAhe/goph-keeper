#!/usr/bin/env bash

#go run -ldflags "-X main.Version=v1.0.1 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'" main.go

# build with debug info
go build \
    -ldflags "\
        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.Version=1.0.0' \
        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')' \
        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.Stage=DEV' \
        " \
    -o ./cmd/server \
    ./cmd/server/.

# build without debug info
#go build \
#    -ldflags "-s -w \
#        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.Version=1.0.0' \
#        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')' \
#        -X 'github.com/ElfAstAhe/goph-keeper/internal/app/config.Stage=PROD' \
#        " \
#    -o ./cmd/server \
#    ./cmd/server/.