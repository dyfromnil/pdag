dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

exper:
	sudo rm -rf experiment
	go build -o experiment/client/main cmd/client/main.go
	go build -o experiment/server1/main cmd/server/main.go
	cp -r experiment/server1/ experiment/server2/
	cp -r experiment/server1/ experiment/server3/
	cp -r experiment/server1/ experiment/server4/

dispatchMsp:
	sudo cp -r experiment/server1/Keys experiment/server2/Keys
	sudo cp -r experiment/server1/Keys experiment/server3/Keys
	sudo cp -r experiment/server1/Keys experiment/server4/Keys