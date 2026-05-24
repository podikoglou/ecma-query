.PHONY: build
build: internal/spec/spec.html
	go build ./...

internal/spec/spec.html: ecma262/spec.html
	cp $< $@
