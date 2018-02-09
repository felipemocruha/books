SHELL:= /bin/bash


proto:
	protoc -I book/ book/book.proto --go_out=plugins=grpc:book

server:
	cd server && go build -i -o books

client:
	cd grpc_client && go build -i -o books-cli

dump_packages:
	sudo tcpdump -A -s 0 'port 9000' -i lo

test_server:
	cd server && go test -cover


.PHONY: build proto client dump_packages test_client test_server
