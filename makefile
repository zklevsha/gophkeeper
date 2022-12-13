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

certs:
	mkdir certs && \
	echo "Generating CA cert and key" && \
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout certs/ca-key.pem -out certs/ca-cert.pem --subj "/C=RU/L=Moscow/O=Practicum/OU=Practicum/CN=CA server" && \
	echo "Generate server key and sign request" && \
	openssl req -newkey rsa:4096 -nodes -keyout certs/server-key.pem -out certs/server-req.pem -subj "/C=RU/L=Moscow/O=Practicum/OU=Practicum/CN=CA server" && \
	echo "Generating servers cert" && \
	openssl x509 -req -in certs/server-req.pem -days 120 -CA certs/ca-cert.pem -CAkey certs/ca-key.pem -CAcreateserial -out certs/server-cert.pem