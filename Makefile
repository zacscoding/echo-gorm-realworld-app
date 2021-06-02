.PHONY: docs

docs:
	redoc-cli bundle ./docs/swagger.json -o ./docs/doc.html
