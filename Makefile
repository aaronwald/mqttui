# Makefile for MQTT TUI Browser

.PHONY: build run clean test help

# Default target
all: build

# Build the application
build:
	@echo "Building MQTT TUI Browser..."
	go build -o mqttui

# Run with default settings
run: build
	@echo "Starting MQTT TUI with default broker (localhost:1883)..."
	./mqttui

# Run with public test broker
run-test: build
	@echo "Starting MQTT TUI with HiveMQ public broker..."
	MQTT_BROKER="tcp://broker.hivemq.com:1883" ./mqttui

# Run with Eclipse test broker
run-eclipse: build
	@echo "Starting MQTT TUI with Eclipse public broker..."
	MQTT_BROKER="tcp://mqtt.eclipseprojects.io:1883" ./mqttui

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f mqttui

# Run Go tests
test:
	@echo "Running tests..."
	go test -v ./...

# Vet the code
vet:
	@echo "Running go vet..."
	go vet ./...

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o mqttui-linux
	GOOS=darwin GOARCH=amd64 go build -o mqttui-macos
	GOOS=windows GOARCH=amd64 go build -o mqttui.exe

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Show help
help:
	@echo "MQTT TUI Browser - Available commands:"
	@echo ""
	@echo "  build       Build the application"
	@echo "  run         Run with default broker (localhost:1883)"
	@echo "  run-test    Run with HiveMQ public broker"
	@echo "  run-eclipse Run with Eclipse public broker"
	@echo "  clean       Clean build artifacts"
	@echo "  test        Run tests"
	@echo "  vet         Run go vet"
	@echo "  fmt         Format code"
	@echo "  build-all   Build for multiple platforms"
	@echo "  deps        Install/update dependencies"
	@echo "  help        Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  MQTT_BROKER   - MQTT broker URL (default: tcp://localhost:1883)"
	@echo "  MQTT_USERNAME - MQTT username (optional)"
	@echo "  MQTT_PASSWORD - MQTT password (optional)"
	@echo "  MQTT_CLIENT_ID - MQTT client ID (default: mqttui)"