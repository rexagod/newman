export GO111MODULE=on

# Install the required dependencies.
install:
	@go mod tidy

run:install
	go run .