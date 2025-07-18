.PHONY: run/main build/linux build/mac build/win

LINUX_BINARY_DIR="bin/linux"
MAC_BINARY_DIR="bin/mac"
WIN_BINARY_DIR="bin/win"

run/main:
	@go run ./cmd/app

build/linux:
	@mkdir -p $(LINUX_BINARY_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(LINUX_BINARY_DIR)/pomodoro-cli ./cmd/app

build/mac:
	@mkdir -p $(MAC_BINARY_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(MAC_BINARY_DIR)/pomodoro-cli ./cmd/app

build/win:
	@mkdir -p $(WIN_BINARY_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(WIN_BINARY_DIR)/pomodoro-cli.exe ./cmd/app
