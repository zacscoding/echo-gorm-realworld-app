.PHONY: docs generate test cleantests

cleantests:
	go clean -testcache

test:
	go test ./...

generate:
	go generate ./...

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
