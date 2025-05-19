ARG GO_VERSION="1.22"
ARG RUNNER_IMAGE="gcr.io/distroless/static-debian11"
ARG BUILD_TAGS="netgo,ledger,muslc"
ARG GIT_VERSION="mainnet-v2.0.0"
ARG GIT_COMMIT=""

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine as builder

ARG BUILD_TAGS

RUN apk add --no-cache \
    ca-certificates \
    build-base \
    linux-headers \
    binutils-gold \
    git

WORKDIR /app

# Copy go.mod and go.sum
COPY ../go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY ../.. .

# Cosmwasm - Download ARM64 libwasmvm
RUN wget https://github.com/CosmWasm/wasmvm/releases/download/v1.5.4/libwasmvm_muslc.aarch64.a \
    -O /lib/libwasmvm_muslc.a

# Build selfchain binary for ARM64
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
    go build \
    -mod=readonly \
    -tags "${BUILD_TAGS}" \
    -ldflags \
    "-X github.com/cosmos/cosmos-sdk/version.Name="selfchain" \
    -X github.com/cosmos/cosmos-sdk/version.AppName="selfchaind" \
    -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static' \
    -X github.com/cosmos/cosmos-sdk/version.Version=${GIT_VERSION} \
         -X github.com/cosmos/cosmos-sdk/version.Commit=${GIT_COMMIT}" \
    -trimpath \
    -o /app/build/selfchaind-linux-arm64 \
    ./cmd/selfchaind/main.go

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM alpine:latest
COPY --from=builder /app/build/selfchaind-linux-arm64 /usr/local/bin/selfchaind
# Add any additional configuration for the runner stage