
build:
	go build -o minecraftd -trimpath -ldflags=-buildid=

test:
	go test

fmt:
	go fmt
