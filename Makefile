.PHONY=run
run:
	go build && ./TaigaBot

.PHONY=build
build:
	go build -o main