APP_NAME=donut

.PHONY: build clean

build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

clean:
	rm -f bin/$(APP_NAME)
