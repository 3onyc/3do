SHELL=/bin/bash

BINARY=3do
PACKAGE=github.com/3onyc/3do
BOWER_BIN=node_modules/bower/bin/bower

dist: 3do frontend-dist deps
	rice append --exec=$(BINARY)

frontend-dist:
	cd frontend && ember build -prod

3do: deps
	bunch go build -o $(BINARY) $(PACKAGE)

deps: frontend-deps backend-deps

frontend-deps: frontend/$(BOWER_BIN)
	cd frontend && npm install
	cd frontend && $(BOWER_BIN) install

frontend/$(BOWER_BIN):
	cd frontend && npm install bower

backend-deps: $(GOPATH)/bin/bunch $(GOPATH)/bin/rice
	bunch install

$(GOPATH)/bin/bunch:
	go get github.com/dkulchenko/bunch

$(GOPATH)/bin/rice:
	go get github.com/GeertJohan/go.rice/rice

debug: backend-deps
	bunch go build -o $(BINARY) $(PACKAGE)

test: backend-deps
	GOPATH=$(PWD)/.vendor go get -v github.com/google/gofuzz
	
	shopt -s nullglob; \
	for F in $$(find . -type d -not -path '*/\.*' -a -not -path "*frontend*"); do \
		if [ -n "$$(echo $$F/*_test.go)" ]; then \
			bunch go test "$(PACKAGE)/$${F:2}"; \
		fi; \
	done

clean:
	rm $(BINARY)

.PHONY: dist frontend-dist deps frontend-deps backend-deps debug test clean
