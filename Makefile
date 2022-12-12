run: 
	go run cmd/main.go
migrateup:
	migrate -path ./migrations/postgres -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable' up
migratedown:
	migrate -path ./migrations/postgres -database 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable' down
pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge



	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/auth_db?sslmode=disable' up
