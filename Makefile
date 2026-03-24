service_run:
	go run cmd/sso/main.go


docker_compose_up:
	docker-compose up -d
docker_compose_down:
	docker-compose down