build:
	go build -o ./cmd/simplai/simplai ./cmd/simplai/...

run: build
	@ ./cmd/simplai/simplai

default: build
