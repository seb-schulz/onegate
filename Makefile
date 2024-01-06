UI_BASEDIR := internal/ui

.PHONY: build
build: build-client
	CGO_ENABLED=0 go build -v -ldflags '-w -extldflags "-static"' -tags embedded,production ./

.PHONY: build-client
build-client:
	go generate -tags embedded,production ./...

.PHONY: generate
generate:
	go generate -v -tags gqlgen ./...

.PHONY: test vet
test vet:
	go $@ -tags embedded,production ./...

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
