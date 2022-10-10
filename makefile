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

