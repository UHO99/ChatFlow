.PHONY: createdb dropdb migrateup migratedown sqlc

DB_URL=postgres://root:1234@localhost:5432/test?sslmode=disable

createdb:
	docker exec -it test_postgres createdb -U postgres test_postgres

dropdb:
	docker exec -it test_postgres dropdb -U postgres test_postgres

migrateup:
	migrate -path GoChat/db/migration -database $(DB_URL) up

migratedown:
	migrate -path GoChat/db/migration -database $(DB_URL) down

sqlc:
	sqlc generate -f GoChat/sqlc.yaml