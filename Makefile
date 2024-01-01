UI_BASEDIR := internal/ui

.PHONY: build
build: generate
	CGO_ENABLED=0 go build -v -ldflags '-w -extldflags "-static"' -tags embedded,production ./

.PHONY: generate
generate:
	go generate -v -tags embedded,production ./...

.PHONY: migrate
migrate:
	go run -v -tags embedded,production ./ migrate

.PHONY: clean
clean: clean-ui clean-go

.PHONY: clean-ui
clean-ui:
	@rm -rf internal/ui/_build/*

.PHONY: clean-go
clean-go:
	@rm -f onegate
