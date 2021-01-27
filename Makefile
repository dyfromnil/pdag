dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

run:
	rm -rf store
	go run cmd/main.go