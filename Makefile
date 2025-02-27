ifndef SHA
	SHA := $(shell git rev-parse HEAD)
endif

ifndef BUILD_TIME
	BUILD_TIME := $(shell git show -s --format=%ct $(SHA))
endif

GOMOD := $(shell go list -m)
GOBUILD_FLAGS := -ldflags "-X $(GOMOD)/internal/build.vcsInfo=$(SHA),$(BUILD_TIME)"

.PHONY: all clean publish

ALL: bin/reader

bin/%: cmd/%/main.go $(shell find pkg -name '*.go')
	go build -o $@ $(GOBUILD_FLAGS) ./cmd/$*

bin/buildimg:
	go build -o $@ github.com/kellegous/buildimg

reader.tar: Dockerfile $(shell find cmd pkg -type f) bin/buildimg
	bin/buildimg --tag=$(TAG) --target=linux/amd64:$@ --build-arg=SHA=${SHA} --build-arg=BUILD_TIME=${BUILD_TIME} kellegous/reader

publish: reader.tar
	sup host image load @ $<

clean:
	rm -rf bin