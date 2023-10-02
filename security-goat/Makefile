.PHONY: all build-linux build-macos build-windows generate-checksum clean

BINARY_NAME=security-goat
RELEASE_PATH=../build

release: build-linux build-macos build-windows generate-checksum

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-linux-amd64 main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-macos-amd64 main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(RELEASE_PATH)/$(BINARY_NAME)-windows-amd64.exe main.go

generate-checksum:
	cd $(RELEASE_PATH) && \
	shasum -a 256 $(BINARY_NAME)-linux-amd64 > $(BINARY_NAME)-linux-amd64.sha256 && \
	shasum -a 256 $(BINARY_NAME)-macos-amd64 > $(BINARY_NAME)-macos-amd64.sha256 && \
	shasum -a 256 $(BINARY_NAME)-windows-amd64.exe > $(BINARY_NAME)-windows-amd64.sha256

clean:
	rm -f $(RELEASE_PATH)/$(BINARY_NAME)-* $(RELEASE_PATH)/*.sha256

run:
	go run main.go

# Make the module available by running the go list command to prompt Go to update its index of modules with information about the module you’re publishing.
go-publish:
	GOPROXY=proxy.golang.org go list -m github.com/andersonbosa/security-goat@latest
	go install github.com/andersonbosa/security-goat@latest

docker:
	docker build -t security-goat:latest . && docker run --rm -it security-goat:latest 

