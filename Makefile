runc:
	docker run --rm -d -p 8001:80 --name=nginx1 nginx 
	docker run --rm -d -p 8002:80 --name=nginx2 nginx
	docker run --rm -d -p 8003:80 --name=nginx3 nginx

stop:
	docker stop nginx1 nginx2 nginx3

run-server: cmd/server/main.go
	go run cmd/server/main.go

run-all: runc run-server

.PHONY: runc stop run-server run-all
