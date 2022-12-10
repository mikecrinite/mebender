dockerize:
	docker build --tag mebender .

run-docker: dockerize
	docker-compose up

lint:
	go fmt ./...
	go vet ./...