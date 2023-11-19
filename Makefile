UI_BASEDIR := internal/ui

.PHONY: build
build: generate
	go build -v -tags embedded,production ./cmd/onegate

.PHONY: generate
generate:
	go generate -v -tags embedded,production ./...

.PHONY: clean
clean: clean-ui clean-go

.PHONY: clean-ui
clean-ui:
	@rm -rf internal/ui/_build/*

.PHONY: clean-go
clean-go:
	@rm -f onegate
