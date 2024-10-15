gen_postgres_migration:
	migrate create -ext=sql -dir=./postgres/migration -seq $(name)

docker-build:
	docker build -t xybor/todennus-migration -f ./Dockerfile .

docker-compose-up:
	docker compose --env-file .env -f ./compose.yaml up -d

docker-compose-down:
	docker compose -f ./compose.yaml down
