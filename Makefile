BINARY=3do
PACKAGE=github.com/3onyc/3do

dist:
	cd frontend && ember build -prod
	go build -o $(BINARY) $(PACKAGE)
	rice append --exec=$(BINARY)

deps: frontend-deps backend-deps

frontend-deps:
	cd frontend && npm install
	cd frontend && bower install

backend-deps:
	go install -v $(PACKAGE)/...

debug:
	go build -o $(BINARY) $(PACKAGE)

clean:
	rm $(BINARY)

.PHONY: dist clean deps frontend-deps backend-deps
