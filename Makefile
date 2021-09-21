dag:
	rm -rf proto-go/common/*
	protoc --go_out=plugins=grpc:. ./blockproto/block.proto	

build:export CGO_ENABLED=0
build:
	@echo $(CGO_ENABLED)
	sudo rm -rf experiment
	go build -o experiment/client/main cmd/client/main.go
	go build -o experiment/server0/main cmd/server/main.go
	go build -o experiment/genKeys cmd/genKeys/main.go

	cd experiment/ && ./genKeys
	mv experiment/Keys experiment/server0/
	
	cp -r experiment/server0/ experiment/server1/
	cp -r experiment/server0/ experiment/server2/
	cp -r experiment/server0/ experiment/server3/
	# cp -r experiment/server0/ experiment/server4/
	# cp -r experiment/server0/ experiment/server5/
	# cp -r experiment/server0/ experiment/server6/
	# cp -r experiment/server0/ experiment/server7/
	# cp -r experiment/server0/ experiment/server8/
	# cp -r experiment/server0/ experiment/server9/

cli:
	go build -o experiment/client/main cmd/client/main.go

export DELAY=40ms
export JITTER=10ms
# export RATE=2mbit
start_server_container:
	cd docker/ && ./container.sh stopServer 
	cd docker/ && ./container.sh constructServerImage && DELAY=$(DELAY) JITTER=$(JITTER) RATE=$(RATE) docker-compose -f docker-compose-servers.yaml up
	# iftop -nN -i eth0

start_client_container:
	cd docker/ && ./container.sh stopClient 
	cd docker/ && ./container.sh constructClientImage
	cd docker/ && docker-compose -f docker-compose-client.yaml up -d
	docker exec -it client /bin/sh

stop_all:
	cd docker/ && ./container.sh stopClient
	cd docker/ && ./container.sh stopServer

stop_client:
	cd docker/ && ./container.sh stopClient

stop_server:
	cd docker/ && ./container.sh stopServer

base_image:
	cd docker/ && ./container.sh constructBaseImage

add_net_delay:
	# pumba netem -d 10m --tc-image gaiadocker/iproute2 rate --rate 100kbit re2:server*
	# pumba netem -d 10m --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal rate --rate 100kbit re2:server*
	# pumba netem -d 10m --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal re2:server*
	# pumba netem -d 10m --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal server1 server2 server3
	# pumba netem -d 10m --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal re2:server*
	pumba netem -d 1h --tc-image gaiadocker/iproute2 delay --time 40 --jitter 10 --distribution normal rate --rate 100kbit server0 server1 server2 server3
