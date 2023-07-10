SHA := $(shell git rev-parse HEAD)
TAG := $(shell git rev-parse --short HEAD)

ALL: bin/reader

bin/%: cmd/%/main.go $(shell find pkg -name '*.go')
	go build -o $@ ./cmd/$*

bin/buildimg:
	go build -o $@ github.com/kellegous/buildimg

reader-$(TAG).tar: Dockerfile $(shell find cmd pkg -type f) bin/buildimg
	bin/buildimg --tag=$(TAG) --target=linux/amd64:$@ kellegous/reader

publish: reader-$(TAG).tar

clean:
	rm -rf bin