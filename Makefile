.PHONY: docs generate test lint cleantests compose

cleantests:
	@go clean -testcache

test:
	@go test ./...

lint:
	@gofmt -d .

generate:
	@go generate ./...

compose.%:
	$(eval CMD=${subst compose.,,$(@)})
	./scripts/compose.sh ${CMD}

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
