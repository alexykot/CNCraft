# Any "makefile:XX: *** missing separator.  Stop." means you used spaces instead of tabs.
SHELL := /bin/bash
# magic
# (catch-all rule to avoid make throwing "nonexistent rule" errors for parameters passed into existing rules)
%:
	@true

# more magic
# (avoid clash of `build` command with `build` directory)
.PHONY: build

format:
	@go fmt ./pkg/... ./cmd/... ./core/...
	@goimports -format-only -local "github.com/alexykot/cncraft" -l -w ./pkg ./cmd ./core

run:
	@reset
	@go run ./cmd/server/server.go

gen:
	cmd/tools/build/generate.sh $(filter-out $@, $(MAKECMDGOALS))

test:
	@go test -p 1 --count 1 ./...

run-client:
	@reset
	@go run ./cmd/client/client.go

idkfa:
	@go run ./cmd/tools/main.go misc idkfa Kolsar
