# GOOS = darwin
# GOARCH = arm64
GOOS = linux
GOARCH = amd64

PROJECT_NAME = keep-notes
MAIN_FILE = cmd/server/main.go
BUILD_PATH = build
ENV_FILE = oracle.env
SHELL := /bin/bash


init-build-dirs:
	@rm -rf $(BUILD_PATH)
	@mkdir $(BUILD_PATH)
	$(info $(BUILD_PATH) was created)
build-api: init-build-dirs
	$(info Will build API for $(GOOS) and $(GOARCH))
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(BUILD_PATH)/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "API build completed"
clean:
	@rm -rf $(BUILD_PATH)
deploy: build-api
	scp $(BUILD_PATH)/$(PROJECT_NAME) skrao@api.skrao.net:~/Dev/apis/keep-notes/
run:
	go run $(MAIN_FILE)
lint:
	$(HOME)/go/bin/golangci-lint run