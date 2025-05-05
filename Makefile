docker-rebuild:
	docker-compose -f docker-compose.yml --env-file .env stop
	DOCKER_BUILDKIT=1 docker-compose -f docker-compose.yml --env-file .env build
	docker-compose -f docker-compose.yml --env-file .env up -d --remove-orphans

docker-rebuild-local:
	docker compose -f docker-compose-local.yml stop
	docker compose -f docker-compose-local.yml build
	docker compose -f docker-compose-local.yml up -d --remove-orphans

docker-rebuild-test:
	docker compose -f docker-compose-test.yml --env-file .env.test stop
	docker compose -f docker-compose-test.yml --env-file .env.test build
	docker compose -f docker-compose-test.yml --env-file .env.test up -d --remove-orphans
