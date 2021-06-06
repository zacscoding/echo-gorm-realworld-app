.PHONY: docs generate

generate:
	go generate ./...

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
