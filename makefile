server:
	cd ./cmd/server/ && go build .

client:
	cd ./cmd/client/ && go build .

pb :
	protoc --proto_path=proto proto/*.proto --go_out=internal --go-grpc_out=internal

lint:
	go vet ./... 
	staticcheck ./...  
	errcheck ./... 
	golint ./...

db:
	psql -U gophkeeper -d  gophkeeper  -f sql/schema.sql && \
	psql -U gophkeeper -d  gophkeeper  -f sql/data.sql

drop:
	psql -U gophkeeper -d  gophkeeper  -f sql/drop.sql

# dont forget to set env variable
# export POSTGRESQL_URL='postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable'
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down