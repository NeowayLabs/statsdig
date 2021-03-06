version=latest

build:
	docker build . -t neowaylabs/statsdig:$(version)

publish: build
	docker push neowaylabs/statsdig:$(version)

check: 
	go test -race -timeout 10s ./...
