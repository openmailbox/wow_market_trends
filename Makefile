.PHONY: clean deps go js

EXE = ./cmd/wowexchange/wowexchange
WEBDIR = ./web/static/src
JSDIR = $(WEBDIR)/scripts
CSSDIR = $(WEBDIR)/styles

all: deps go js

js: $(wildcard ./web/static/src/scripts/*.js)
	npm run build

go: $(wildcard ./internal/*.go) $(wildcard ./cmd/wowexchange/*.go)
	go build ./internal
	go build -o $(EXE) ./cmd/wowexchange

deps:
	npm install --loglevel=error

clean: 
	rm $(EXE)
	rm web/static/dist/*.js
