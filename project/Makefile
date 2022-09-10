BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=logAPP
LISTENER_BINARY=listenerService
up:
	@echo "Starting docker image"
	docker-compose up -d
	@echo "Docker image started"


up_build: build_broker build_auth build_logger build_listener
	@echo "Stop docker images (if running...)"
	docker-compose down
	@echo "Build (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker image build and started..."


down:
	@echo "Stop docker compose..."
	docker-compose down
	@echo "Done!"

build_broker:
	@echo "build borker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

clean:
	@echo "Remove dangling images..."
	docker image prune
	@echo "Done!"

build_auth:
	@echo "build auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

build_logger:
	@echo "Build logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

build_listener:
	@echo "Build listener service binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

