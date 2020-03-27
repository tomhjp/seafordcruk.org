all: test handler

handler: main.go
	GOOS=linux go build -o handler
	zip handler.zip handler

test: main.go main_test.go
	go test

.PHONY: clean
clean:
	rm handler handler.zip
	touch main.go

