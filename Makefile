APP_NAME=donut
COVERAGE_OUT=cover.out
COVERAGE_HTML=cover.html
PACKAGE_NAME=github.com/gleamsoda/donut

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

.PHONY: clean
clean:
	rm -f bin/$(APP_NAME)

.PHONY: test
test:
	go test -cover ./... -coverprofile=$(COVERAGE_OUT)

.PHONY: cover
cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	open $(COVERAGE_HTML)
