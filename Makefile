dockerize:
	docker build --tag mebender .

run-docker: dockerize
	docker-compose up

shell-docker: 
	docker exec -it mebender /bin/sh

lint:
	go fmt ./...
	go vet ./...
