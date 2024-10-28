default:
    @just --list

# build nats-readiness binary
build:
    go build -o nats-readiness ./cmd/nats-readiness

# update go packages
update:
    @cd ./cmd/nats-readiness && go get -u

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
