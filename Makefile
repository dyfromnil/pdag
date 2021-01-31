dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

exper:
	sudo rm -rf experiment
	go build -o experiment/client/main cmd/client/main.go
	go build -o experiment/server0/main cmd/server/main.go
	go build -o experiment/genKeys cmd/server/main.go

	cd experiment/ && ./genKeys
	cp -r experiment/Keys experiment/server0/
	
	cp -r experiment/server0/ experiment/server1/
	cp -r experiment/server0/ experiment/server2/
	cp -r experiment/server0/ experiment/server3/

client:
	go build -o experiment/client/main cmd/client/main.go