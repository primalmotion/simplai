default: build

build:
	go build ./...
	cd ./cmd/simplai && go build

lint:
	golangci-lint run \
		--timeout=5m \
		--disable-all \
		--exclude-use-default=false \
		--exclude=package-comments \
		--exclude=unused-parameter \
		--enable=errcheck \
		--enable=goimports \
		--enable=ineffassign \
		--enable=revive \
		--enable=unused \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=misspell \
		--enable=prealloc \
		--enable=nakedret \
		--enable=typecheck \
		--enable=unparam \
		--enable=gosimple \
		--enable=nilerr \
		./...

test:
	go test ./... -race -cover -covermode=atomic -coverprofile=unit_coverage.out

sec:
	gosec -quiet ./...

vuln:
	govulncheck ./...
