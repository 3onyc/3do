BINARY=3do
PACKAGE=github.com/3onyc/3do

dist:
	cd frontend && ember build -prod
	gopm build -o $(BINARY)
	rice append --exec=$(BINARY)

deps: frontend-deps backend-deps

frontend-deps:
	cd frontend && npm install
	cd frontend && bower install

backend-deps:
	gopm get

debug:
	gopm build -o $(BINARY)

clean:
	rm $(BINARY)

.PHONY: dist clean deps frontend-deps backend-deps
