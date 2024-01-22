bin/bootstrap.zip: bin/bootstrap
	cd bin; zip bootstrap bootstrap

bin/bootstrap: *.go
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap

fmt:
	go fmt

clean:
	rm bin/*