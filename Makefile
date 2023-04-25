VERSION := 0.1.3

devup:
	docker-compose -f scripts/docker-compose.yaml -p api_dev --env-file scripts/dev.env up -d

devdown:
	docker-compose -f scripts/docker-compose.yaml -p api_dev down -v

build:
	docker build --build-arg $(VERSION) -t solte/op-bot:$(VERSION) .

ci_up:
	docker compose -f ./deployments/docker-compose.yaml up -d

ci_teardown:
	docker compose -f ./deployments/docker-compose.yaml down -v