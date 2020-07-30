PROJECT_NAME = go-authorizer
MODULE_NAME = cmd

.SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help:
	$(info go-authorizer commands:)
	$(info -> setup                   install dependencies)
	$(info -> format                  formats go files)
	$(info -> build                   build binary)
	$(info -> test                    runs available tests)
	$(info -> run                     runs application)
	$(info -> docker-build            builds application on a docker image)
	$(info -> docker-test             runs available tests on a docker image)
	$(info -> docker-run              runs application on a docker image)

.PHONY: setup
install:
	go get -d -v -t ./...
	go install -v ./...
	go mod tidy -v

.PHONY: format
format:
	go fmt ./...

.PHONY: build
build:
	go build -v -o $(MODULE_NAME).bin ./$(MODULE_NAME)
	chmod +x $(MODULE_NAME).bin
	echo $(MODULE_NAME).bin

.PHONY: test
test:
	go test ./... -v -covermode=count

.PHONY: run
run:
	go run ./$(MODULE_NAME)

.PHONY: docker-build
docker-build:
	docker build --build-arg root_dir=./$(MODULE_NAME) -t $(PROJECT_NAME) .

.PHONY: docker-run
docker-run: docker-build
	docker run -a stdin -a stdout -i -t --name $(PROJECT_NAME) --rm $(PROJECT_NAME)

.PHONY: docker-test
docker-test: docker-build
	docker run --name $(PROJECT_NAME) --rm $(PROJECT_NAME) go test ./... -v -covermode=count
