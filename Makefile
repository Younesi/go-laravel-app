BINARY_NAME=atlasApp

build:
	@go mod tidy
	@go mod vendor
	@echo "Building Atlas..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Atlas built!"

run: build
	@echo "Starting Atlas..."
	@./tmp/${BINARY_NAME} &
	@echo "Atlas started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping Atlas..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Atlas!"

restart: stop start