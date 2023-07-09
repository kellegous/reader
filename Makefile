ALL: bin/reader

bin/%: cmd/%/main.go $(shell find pkg -name '*.go')
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin