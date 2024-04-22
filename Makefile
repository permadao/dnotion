
build_local:
	go build -o ./build/dnotion ./run/service

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/dnotion_linux ./run/service

TARGET=dnotion
VERSION=0.0.1
# Makefile for building docker image
DOCKER_TAG=$(TARGET):$(VERSION)
DOCKERFILE_PATH=Dockerfile

docker:
	docker build -t $(DOCKER_TAG) -f $(DOCKERFILE_PATH) .

docker_run:
	@echo "Running docker container with tag: $(DOCKER_TAG)"
	docker run -d --name $(TARGET)  $(DOCKER_TAG)