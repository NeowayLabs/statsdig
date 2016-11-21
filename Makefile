version=latest

build:
	docker build . -t statsdig:$(version)
