dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

exper:
	sudo rm -rf experiment
	go build -o experiment/client/main cmd/client/main.go
	go build -o experiment/server0/main cmd/server/main.go
	cp -r experiment/server0/ experiment/server1/
	cp -r experiment/server0/ experiment/server2/
	cp -r experiment/server0/ experiment/server3/

client:
	go build -o experiment/client/main cmd/client/main.go

dispatchMsp:
	sudo cp -r experiment/server0/Keys experiment/server1/Keys
	sudo cp -r experiment/server0/Keys experiment/server2/Keys
	sudo cp -r experiment/server0/Keys experiment/server3/Keys