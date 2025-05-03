.PHONY: protogen

protogen: 
	protoc -Ipb --go_out=. --go_opt=module=github.com/zero-shubham/authsvc --go-grpc_out=. --go-grpc_opt=module=github.com/zero-shubham/authsvc pb/service.proto

get_migrate:
	go get -u -d github.com/golang-migrate/migrate/cmd/migrate

create_migrate: 
	migrate create -ext sql -dir db/migrations/ -seq $(seq)

migrate_up:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/datab?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/datab?sslmode=disable" -verbose down

protobin:
	protoc --descriptor_set_out=pb/service.pb pb/service.proto

seed:
	go run ./cmd/seed/main.go

test:
	export ENV=test && docker compose -f docker-compose-test.yaml run authsvc-test go test ./... -v
	docker compose -f docker-compose-test.yaml down --remove-orphans