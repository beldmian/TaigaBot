.PHONY=run
run:
	docker run --env-file ./env taiga_bot:latest 

.PHONY=build
build:
	docker build . -t taiga_bot:latest