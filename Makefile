build-all:
	GOOS=linux GOARCH=amd64 go build -tags ledger -o ./build/selfchaind-linux-amd64 ./cmd/selfchaind/main.go
	GOOS=linux GOARCH=arm64 go build -tags ledger -o ./build/selfchaind-linux-arm64 ./cmd/selfchaind/main.go
	GOOS=darwin GOARCH=amd64 go build -tags ledger -o ./build/selfchaind-darwin-amd64 ./cmd/selfchaind/main.go
	GOOS=darwin GOARCH=arm64 go build -tags ledger -o ./build/selfchaind-darwin-arm64 ./cmd/selfchaind/main.go

build-cosmovisor:
	go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
	mkdir -p ./build
	mv /go/bin/cosmovisor ./build/cosmovisor

do-checksum:
	cd build && sha256sum \
		selfchaind-linux-amd64 selfchaind-linux-arm64 \
		selfchaind-darwin-amd64 selfchaind-darwin-arm64 \
		> selfchain_checksum

build-with-checksum: build-all do-checksum
