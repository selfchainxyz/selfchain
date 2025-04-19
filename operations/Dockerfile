ARG GO_VERSION="1.22"
ARG RUNNER_IMAGE="gcr.io/distroless/static-debian11"
ARG BUILD_TAGS="netgo,ledger,muslc"

ARG COSMOVISOR_VERSION="v1.5.0"
ARG github_token
ARG version

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:${GO_VERSION}-alpine3.20 as builder

ARG GIT_VERSION
ARG GIT_COMMIT
ARG BUILD_TAGS
ARG COSMOVISOR_VERSION
ARG github_token
ARG version

RUN apk add --no-cache \
    ca-certificates \
    build-base \
    linux-headers \
    binutils-gold \
    git

# Clone repo
RUN git clone https://$github_token:x-oauth-basic@github.com/selfchainxyz/selfchain.git /selfchain

# Download go dependencies
WORKDIR /selfchain
RUN git checkout $version
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Cosmwasm - Download correct libwasmvm version
RUN ARCH=$(uname -m) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/v1.5.4/libwasmvm_muslc.$ARCH.a \
    -O /lib/libwasmvm_muslc.$ARCH.a && \
    cp /lib/libwasmvm_muslc.$ARCH.a /lib/libwasmvm_muslc.a

# Build selfchain binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    GOWORK=off go build \
    -mod=readonly \
    -tags "netgo,ledger,muslc" \
    -ldflags \
    "-X github.com/cosmos/cosmos-sdk/version.Name="selfchain" \
    -X github.com/cosmos/cosmos-sdk/version.AppName="selfchaind" \
    -X github.com/cosmos/cosmos-sdk/version.Version=${GIT_VERSION} \
    -X github.com/cosmos/cosmos-sdk/version.Commit=${GIT_COMMIT} \
    -X github.com/cosmos/cosmos-sdk/version.BuildTags=${BUILD_TAGS} \
    -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
    -trimpath \
    -o ./build/selfchaind \
    ./cmd/selfchaind/main.go

# --------------------------------------------------------
# Cosmovisor runner
# --------------------------------------------------------

FROM ${RUNNER_IMAGE}

COPY --from=builder /selfchain/build/selfchaind /bin/selfchaind

ENV HOME=/selfchain
WORKDIR $HOME

EXPOSE 26656
EXPOSE 26657
EXPOSE 1317

ENTRYPOINT ["selfchaind"]
