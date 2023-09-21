#TARGET
TARGET=dnotion
VERSION=0.0.1

PORT=3002

# Makefile for building docker image
DOCKER_TAG=$(TARGET):$(VERSION)

# build the image
DOCKERFILE_PATH=Dockerfile

build:
	@echo "Building: $(TARGET)"
	go build -o $(TARGET) ./start/main.go

docker:
	@echo "Building docker image with tag: $(DOCKER_TAG)"
	docker build -t $(DOCKER_TAG) -f $(DOCKERFILE_PATH) .

run:
	@echo "Running docker container with tag: $(DOCKER_TAG)"
	docker run -d -p 3002:$(PORT) --name $(TARGET)  $(DOCKER_TAG)

stop:
	@echo "Stopping $(TARGET) image..."
	docker stop $(TARGET)

clean:
	@echo "Cleaning up images with tag: $(DOCKER_TAG)"
	docker stop $(TARGET)
	docker rm $(TARGET)
	docker image rm $(DOCKER_TAG)