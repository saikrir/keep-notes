GOOS = darwin
GOARCH = arm64
PROJECT_NAME = keep-notes
MAIN_FILE = cmd/server/main.go
BUILD_PATH = build


init-build-dirs:
	@rm -rf $(BUILD_PATH)
	@mkdir $(BUILD_PATH)
	$(info $(BUILD_PATH) was created)
build-api: init-build-dirs
	$(info Will build API for $(GOOS) and $(GOARCH))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_PATH)/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "API build completed"
cleanup:
	@rm -rf $(BUILD_PATH)
