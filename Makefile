default: test

.PHONY: test
test: start-emulator
	go test

.PHONY: start-emulator
start-emulator:
	docker-compose up -d

.PHONY: stop-emulator
stop-emulator:
	docker-compose down
