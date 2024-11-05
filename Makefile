APP_NAME = proxy
SRC_DIR = ./$(APP_NAME)
BUILD_DIR = ./bin
# TEST_DIR = ./test
GO_FILES = $(shell find . -name '*.go')

.PHONY: all
all: build

.PHONY: build
build: $(GO_FILES)
	@echo "Building..."

	@mkdir -p $(BUILD_DIR)

	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

.PHONY: test
test:
	@echo "Testing..."

	go test ./... -v

.PHONY: clean
clean:
	@echo "Cleaning..."

	@rm -rf $(BUILD_DIR)

.PHONY: lint
lint:
	@echo "Linting..."

	golangci-lint run

.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."

	@$(BUILD_DIR)/$(APP_NAME)
