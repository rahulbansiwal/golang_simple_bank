DB_URL=postgresql://admin:admin@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres-server --network simplebank   -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin  -p 5432:5432 -d  338ccfade89d

createdb:
	docker exec -it postgres-server createdb --username=admin --owner=admin simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down 1

migratedown2:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down 2

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)


dropdb:
	docker exec -it postgres-server dropdb --username=admin  simple_bank

sqlc:
	docker run --rm -v "F:\Code\Go\go-lang_tutorials\backend:/src" -w /src sqlc/sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/mock.go -package mockdb simple_bank/db/sqlc Store

proto:
	 rm -f pb\*.go  
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: createdb postgres dropdb migratedown migrateup sqlc test server mock migratedown1 migrateup1 proto redis new_migration