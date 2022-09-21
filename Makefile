

go ?= GO111MODULE=on go
export go

install:
	go install $(shell $(go) list -m)@$(shell $(MAKE) -f $(firstword $(MAKEFILE_LIST)) go-version)

go-version: go.mod $(shell $(go) list -f '{{$$Dir := .Dir}}{{range .GoFiles}}{{$$Dir}}/{{.}} {{end}}' ./...)
	@TZ=UTC git log -1 '--date=format-local:%Y%m%d%H%M%S' --abbrev=12 '--pretty=tformat:v0.0.0-%cd-%h' $^

go-get:
	@echo $(go) get $(shell $(go) list -m)@$(shell $(MAKE) -f $(firstword $(MAKEFILE_LIST)) go-version)

upgrade-flagx: ../flagx/Makefile
	$(shell $(MAKE) -C ../flagx go-get)
	$(go) mod tidy
upgrade-jsonptr: ../jsonptr/Makefile
	$(shell $(MAKE) -C ../jsonptr go-get)
	$(go) mod tidy
