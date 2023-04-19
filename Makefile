devup:
	docker-compose -f scripts/docker-compose.yaml -p api_dev --env-file scripts/dev.env up -d

devdown:
	docker-compose -f scripts/docker-compose.yaml -p api_dev down -v