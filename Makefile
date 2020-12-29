GO ?= go
VERSION=`git describe --tags`
BINARY=noise

build:
	$(GO) build -o bin/${BINARY}-darwin-amd64 ./cmd/${BINARY}

build_amd64:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-w -s" -o ./bin/${BINARY}-amd64-linux ./cmd/${BINARY}

noise:
	$(GO) build -o bin/${BINARY}-darwin-amd64 ./cmd/${BINARY}
	YOMO_SOURCE_ADDR=localhost:1883 YOMO_ZIPPER_ADDR=localhost:4242 ./bin/${BINARY}-darwin-amd64

build_cli:
	$(GO) build -o bin/yomo-mqtt ./cmd/yomo-mqtt

run_example:
	$(GO) build -o bin/yomo-mqtt ./cmd/yomo-mqtt
	bin/yomo-mqtt run -f example/app.go -p 1883 -z localhost:9999 -t NOISE