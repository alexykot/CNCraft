# Any "makefile:XX: *** missing separator.  Stop." means you used spaces instead of tabs.
SHELL := /bin/bash
# magic
# (catch-all rule to avoid make throwing "nonexistent rule" errors for parameters passed into existing rules)
%:
	@true

# more magic
# (avoid clash of `build` command with `build` directory)
.PHONY: build

lint:
	@go fmt ./pkg/... ./cmd/... ./core/...
	@goimports -format-only -local "github.com/alexykot/cncraft" -l -w ./pkg ./cmd ./core

run:
	go run ./cmd/server/server.go

gen:
	@./proto/gen.sh

