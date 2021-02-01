dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

exper:
	sudo rm -rf experiment
	go build -o experiment/client/main cmd/client/main.go
	go build -o experiment/server0/main cmd/server/main.go
	go build -o experiment/genKeys cmd/genKeys/main.go

	cd experiment/ && sudo ./genKeys
	sudo mv experiment/Keys experiment/server0/
	
	sudo cp -r experiment/server0/ experiment/server1/
	sudo cp -r experiment/server0/ experiment/server2/
	sudo cp -r experiment/server0/ experiment/server3/

cli:
	go build -o experiment/client/main cmd/client/main.go