.PHONY: run docker

run:
	go run ../cmd

docker:
	docker compose up -d