gen_postgres_migration:
	migrate create -ext=sql -dir=./postgres/migration -seq $(name)

migrate:
	go run ./cmd/main.go --postgres

docker-build:
	docker build -t xybor/todennus-migration -f ./Dockerfile .
