build-linux:
	GOOS=linux GOARCH=amd64 go build -o "prop-filter-cli"

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "prop-filter-cli.exe"

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o "prop-filter-cli-mac"

build-all: build-linux build-windows build-mac
	@echo "completed"
