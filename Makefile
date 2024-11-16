build-linux:
	GOOS=linux GOARCH=amd64 go build -o "prop-filter-cli_linux_amd64"

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "prop-filter-cli_windows_amd64.exe"

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o "prop-filter-cli_darwin_amd64"

build-all: build-linux build-windows build-mac
	@echo "completed"
