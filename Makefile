all: build test

build:
	@echo "Generating Windows icon from build/appicon.png..."
	@go run ./tools/genicon -in build/appicon.png -out build/windows/icon.ico
	@echo "Building Wails app..."
	@wails build

run:
	@wails dev

docker-run:
	@docker compose up --build

docker-down:
	@docker compose down

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@powershell -Command "if (Test-Path 'build\\bin') { Remove-Item -Recurse -Force 'build\\bin' }"

.PHONY: all build run test clean
