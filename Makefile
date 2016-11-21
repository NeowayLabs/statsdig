version=latest

build:
	docker build . -t neowaylabs/statsdig:$(version)

push: build
	docker push neowaylabs/statsdig:$(version)
