build:
	@go build -o bin/bin

run: build
	@echo "Starting..."
	@./bin/bin