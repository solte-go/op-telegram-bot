VERSION := 0.1.41

build_worker:
	cd cmd/worker && go build -tags musl -o ../bot -ldflags "-X main.version=$(VERSION)" && cd ../../

docker_build:
	docker build --build-arg $(VERSION) -t solte/op-bot:$(VERSION) .

local_up:
	docker compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose-override.yaml up -d

local_teardown:
	docker compose -f ./deployments/docker-compose.yaml down -v

ci_up:
	docker compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose-test.yaml up -d

ci_teardown:
	docker compose -f ./deployments/docker-compose.yaml down -v