FROM golang:1.20 AS builder

RUN --mount=type=cache,target=/root/.cache \
    apt install git

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG VERSION

RUN  --mount=type=cache,target=/root/.cache \
    cd cmd/worker && go build -tags musl -o ../bot -ldflags "-X main.version=$VERSION" ;


FROM ubuntu:kinetic-20230412 as op-bot
RUN mkdir /etc/worker
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/cmd/prod.toml /opt/prod.toml
COPY --from=builder /build/cmd/bot /opt/worker/

WORKDIR /opt/worker
ENTRYPOINT ["./bot", "-env", "prod"]
