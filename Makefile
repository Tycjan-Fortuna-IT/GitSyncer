.PHONY: dev build lint lint-go lint-ui clean check fmt

dev:
	wails dev

build:
	wails build

lint: lint-go lint-ui

lint-go:
	golangci-lint run ./...

lint-ui:
	cd ui && npm run lint

fmt:
	gofmt -w .
	cd ui && npm run format

check:
	cd ui && npm run check

clean:
	rm -rf build/bin
	rm -rf ui/dist
