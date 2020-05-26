all: test handler

test: main.go main_test.go
	go test

handler: main.go
	GOOS=linux go build -o handler
	zip handler.zip handler

deploy: all
	bash deploy.sh

clean:
	rm handler handler.zip
	touch main.go

.PHONY: all clean test deploy