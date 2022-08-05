.PHONY:fmt
fmt:
	go fmt ./...

.PHONY:vet
vet: fmt
	go vet ./...

.PHONY:build
build: vet
	go build ./...

.PHONY:test_unit
test_unit:
	go test ./...

.PHONY:test_integration
test_integration:
	go test -tags integration ./...

.PHONY: mocks
mocks:
	go generate ./...

.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yml up --build

.PHONY: docker-down
docker-down:
	docker-compose -f docker-compose.yml down
	#docker system prune ## Uncomment this line to remove docker artifacts.
