dag:
	rm -rf proto-go/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto

build:export CGO_ENABLED=0
build:
	@echo $(CGO_ENABLED)
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

start_server_container:
	cd docker/ &&./container.sh && docker-compose -f docker-compose-servers.yaml up

start_client_container:
	# cd docker/ && docker rm client && docker-compose -f docker-compose-client.yaml up
	docker rm client && docker run --name client --network docker_dynet -it client /bin/sh

stop_container:
	cd docker/ &&./container.sh

add_net_delay:
	pumba netem -d 1h --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal re2:server*
