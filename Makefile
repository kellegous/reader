ifndef SHA
	SHA := $(shell git rev-parse HEAD)
endif

ifndef BUILD_TIME
	BUILD_TIME := $(shell git show -s --format=%ct $(SHA))
endif

GOMOD := $(shell go list -m)
GOBUILD_FLAGS := -ldflags "-X $(GOMOD)/internal/build.vcsInfo=$(SHA),$(BUILD_TIME)"

.PHONY: all clean publish reader.tar

ALL: bin/reader

bin/%: cmd/%/main.go $(shell find internal -name '*.go')
	go build -o $@ $(GOBUILD_FLAGS) ./cmd/$*

clean:
	rm -rf bin