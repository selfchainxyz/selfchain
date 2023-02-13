build-all:
	GOOS=linux GOARCH=amd64 go build -o ./build-v2/frontierd-linux-amd64 ./cmd/frontierd/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build-v2/frontierd-linux-arm64 ./cmd/frontierd/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./build-v2/frontierd-darwin-amd64 ./cmd/frontierd/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./build-v2/frontierd-darwin-arm64 ./cmd/frontierd/main.go

build-cosmovisor:
	go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
	mkdir -p ./build-v2
	mv /go/bin/cosmovisor ./build-v2/cosmovisor

do-checksum:
	cd build && sha256sum \
		frontierd-linux-amd64 frontierd-linux-arm64 \
		frontierd-darwin-amd64 frontierd-darwin-arm64 \
		> frontier_checksum

build-with-checksum: build-cosmovisor build-all do-checksum
