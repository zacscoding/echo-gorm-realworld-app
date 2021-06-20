.PHONY: docs generate test lint cleantests

cleantests:
	go clean -testcache

test:
	go test ./...

lint:
	gofmt -d .

generate:
	go generate ./...

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
