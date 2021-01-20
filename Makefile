dag:
	rm -rf proto-go/common
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

run:
	rm -rf store
	go run cmd/main.go