local:
	go run .

deploy: bin/bootstrap.zip
	@AWS_ACCESS_KEY_ID=$(shell cat ./access_key_id) \
	AWS_SECRET_ACCESS_KEY=$(shell cat ./secret_access_key) \
	aws lambda update-function-code \
		--region us-east-1 \
		--function-name lemmeknow-backend \
		--zip-file fileb://bin/bootstrap.zip

invoke:
	@AWS_ACCESS_KEY_ID=$(shell cat ./access_key_id) \
	AWS_SECRET_ACCESS_KEY=$(shell cat ./secret_access_key) \
	aws lambda invoke \
		--region us-east-1 \
		--function-name lemmeknow-backend \
		--cli-binary-format raw-in-base64-out \
		--payload '$(PAYLOAD)' \
		/dev/stderr

bin/bootstrap.zip: bin/bootstrap
	cd bin; zip bootstrap bootstrap

bin/bootstrap: *.go go.*
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap

fmt:
	go fmt

tidy:
	go mod tidy

clean:
	rm bin/*