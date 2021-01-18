dag:
	rm -rf proto-go/common
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto