build-all:
	# Build Darwin arm64 directly
	GOOS=darwin GOARCH=arm64 go build -tags ledger -o ./build/selfchaind-darwin-arm64 ./cmd/selfchaind/main.go

	# Build Linux amd64 using Docker
	docker buildx build --platform linux/amd64 --tag selfchain-builder-amd64:latest -f operations/Dockerfile_amd .
	docker create --name temp-selfchain-container-amd64 selfchain-builder-amd64
	docker cp temp-selfchain-container-amd64:/usr/local/bin/selfchaind ./build/selfchaind-linux-amd64
	docker rm temp-selfchain-container-amd64

	# Build Linux arm64 using Docker
	docker buildx build --platform linux/arm64 --tag selfchain-builder-arm64:latest -f operations/Dockerfile .
	docker create --name temp-selfchain-container-arm64 selfchain-builder-arm64
	docker cp temp-selfchain-container-arm64:/usr/local/bin/selfchaind ./build/selfchaind-linux-arm64
	docker rm temp-selfchain-container-arm64
	
	# Build Darwin amd64 directly
	#GOOS=darwin GOARCH=amd64 go build -tags ledger -o ./build/selfchaind-darwin-amd64 ./cmd/selfchaind/main.go

	# Clean up Docker images
	docker rmi selfchain-builder-amd64:latest selfchain-builder-arm64:latest

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