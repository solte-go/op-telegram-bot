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
    cd cmd/responedr && go build -tags musl -o ../bot -ldflags "-X main.version=$VERSION" \
    cd ../ui && go build -tags musl -o ../op-bot-api -ldflags "-X main.version=$VERSION" ;

#OP-BOT-WORKER
FROM ubuntu:kinetic-20230412 as op-bot-worker
RUN mkdir /opt/worker
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/cmd/prod.toml /opt/prod.toml
COPY --from=builder /build/cmd/bot /opt/worker/

WORKDIR /opt/worker
ENTRYPOINT ["./bot", "-env", "prod"]

#OP-BOT-UI
FROM ubuntu:kinetic-20230412 as op-bot-ui
RUN mkdir /opt/api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/cmd/prod.toml /opt/prod.toml
COPY --from=builder /build/cmd/ui /opt/api/ui/
COPY --from=builder /build/cmd/op-bot-api /opt/api/

WORKDIR /opt/api
ENTRYPOINT ["./op-bot-api", "-env", "prod"]