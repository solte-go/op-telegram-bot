VERSION := 0.1.3

build_worker:
	cd cmd/worker && go build -tags musl -o ../bot -ldflags "-X main.version=$(VERSION)" && cd ../../

docker_build:
	docker build --build-arg $(VERSION) -t solte/op-bot:$(VERSION) .

ci_up:
	docker compose -f ./deployments/docker-compose.yaml up -d

ci_teardown:
	docker compose -f ./deployments/docker-compose.yaml down -v