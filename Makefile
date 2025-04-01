.PHONY: all clean

ALL: bin/reader

bin/%: cmd/%/main.go $(shell find internal -name '*.go')
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin