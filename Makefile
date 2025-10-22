PROTOC_GEN_GO_VERSION := v1.36.10
PROTOC_GEN_TWIRP_VERSION := v8.1.3
PROTOC_VERSION := 33.0

GO_MOD := $(shell go list -m)

BE_PROTOS := \
	reader.pb.go \
	reader.twirp.go

.PHONY: all clean

.PRECIOUS: $(BE_PROTOS)

ALL: bin/reader

bin/%: cmd/%/main.go $(BE_PROTOS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/$*

bin/protoc:
	etc/download-protoc $(PROTOC_VERSION)

bin/protoc-gen-go:
	GOBIN="$(CURDIR)/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

bin/protoc-gen-twirp:
	GOBIN="$(CURDIR)/bin" go install github.com/twitchtv/twirp/protoc-gen-twirp@$(PROTOC_GEN_TWIRP_VERSION)

%.pb.go: %.proto bin/protoc-gen-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go_out=. \
		--go_opt=module=$(GO_MOD) \
		$<

%.twirp.go: %.proto bin/protoc-gen-twirp bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-twirp=bin/protoc-gen-twirp \
		--twirp_out=. \
		--twirp_opt=module=$(GO_MOD) \
		$<

clean:
	rm -rf bin