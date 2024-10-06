init:
	go mod tidy
	go mod vendor
	
start:
	go mod vendor
	docker compose up