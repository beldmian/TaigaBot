.PHONY=run
run:
	go build -o main && ./main

.PHONY=dev
dev:
	mv config.toml prod.toml
	mv dev.toml config.toml

.PHONY=prod
prod:
	mv config.toml dev.toml
	mv prod.toml config.toml
	rm bot.log