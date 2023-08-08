postgres:
	docker run --name postgres-server --network simplebank   -e POSTGRES_PASSWORD=admin -e POSTGRES_USER=admin -p 5432:5432 -d  338ccfade89d

createdb:
	docker exec -it postgres-server createdb --username=admin --owner=admin simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "postgresql://admin:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://admin:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


dropdb:
	docker exec -it postgres-server dropdb --username=admin  simple_bank

sqlc:
	docker run --rm -v "F:\Code\Go\go-lang_tutorials\backend:/src" -w /src sqlc/sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/mock.go -package mockdb simple_bank/db/sqlc Store

.PHONY: createdb postgres dropdb migratedown migrateup sqlc test server mock migratedown1 migrateup1