# Define Go command and flags
GO = go
GOFLAGS = -ldflags="-s -w"

# Define the target executable
TARGET = deviceApi

.PHONY: all clean run test

# Default target: build the executable
all: $(TARGET)

# Rule to build the target executable
$(TARGET): main.go
	$(GO) build $(GOFLAGS) -o $(TARGET) main.go

# Clean target: remove the target executable
clean:
	rm -f $(TARGET)

# Run target: build and run the target executable
run: $(TARGET)
	./$(TARGET)

# Generates the HTML docuementation from the API.yaml file.
docs:
	docker run --rm -i yousan/swagger-yaml-to-html < static/docs/api.yaml > static/docs/api.html

# Test target: run Go tests for the project
test:
	$(GO) test ./...