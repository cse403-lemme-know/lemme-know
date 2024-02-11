.PHONY: dev build test fmt lint tidy clean deploy

dev:
	go run .

build:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap
	cd bin; zip bootstrap bootstrap

test:
	go test

fmt:
	go fmt

tidy:
	go mod tidy

clean:
	rm bin/*

deploy: build
	@AWS_ACCESS_KEY_ID=$(shell cat ./access_key_id) \
	AWS_SECRET_ACCESS_KEY=$(shell cat ./secret_access_key) \
	aws lambda update-function-code \
		--region us-east-1 \
		--function-name lemmeknow-backend \
		--zip-file fileb://bin/bootstrap.zip