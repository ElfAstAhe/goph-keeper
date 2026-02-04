# Переменные для сборки
MODULE_NAME=github.com/ElfAstAhe/goph-keeper
BINARY_NAME=server
BUILD_DIR=./cmd/server
VERSION=1.0.0
BUILD_TIME=$(shell date +'%Y/%m/%d_%H:%M:%S')

.PHONY: build run test clean

# Сборка проекта с прокидыванием переменных
build:
	go build -ldflags "-X '$(MODULE_NAME)/config.Version=$(VERSION)' \
	-X '$(MODULE_NAME)/config.BuildTime=$(BUILD_TIME)'" \
	-o ./bin/$(BINARY_NAME) $(BUILD_DIR)/main.go

# Запуск проекта (сначала соберет, потом запустит)
run: build
	./bin/$(BINARY_NAME)

# Запуск тестов
test:
	go test -v ./...

# Очистка бинарников
clean:
	rm -rf ./bin
