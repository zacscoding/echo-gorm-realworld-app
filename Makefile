.PHONY: docs generate test lint cleantests compose

cleantests:
	@go clean -testcache

tests: test.clean test test.datarace test.build lint

test:
	@echo "Run tests"
	@go test ./... -timeout 5m

test.clean:
	@echo "Clean test cache"
	@go clean -testcache

test.datarace:
	@echo "Run tests with datarace"
	@go test ./... -race -timeout 5m

test.build:
	@echo "Run test build"
	@go build -o /dev/null

lint:
	@echo "Run check lint"
	@gofmt -d .

generate:
	@go generate ./...

compose.%:
	$(eval CMD=${subst compose.,,$(@)})
	./scripts/compose.sh ${CMD}

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
