# Переменные для сборки
PROTO_PATH=api/proto
PROTO_OUT=pkg/api/gophkeeper
MODULE_NAME=github.com/ElfAstAhe/goph-keeper
BINARY_NAME=server
BUILD_DIR=./cmd/server
VERSION=1.0.0
BUILD_TIME=$(shell date +'%Y/%m/%d_%H:%M:%S')

.PHONY: build run test clean

# Генерация gRPC кода
gen-proto:
	mkdir -p $(PROTO_OUT)
	protoc --proto_path=$(PROTO_PATH) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		--go_opt=default_api_level=API_OPAQUE \
		$(PROTO_PATH)/*.proto

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
