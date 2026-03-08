PROTOC_GEN_GO_VERSION := v1.36.10
PROTOC_GEN_CONNECT_GO_VERSION := v1.19.1
PROTOC_VERSION := 33.0

SHA = $(shell go run github.com/kellegous/glue/build/info@latest --format="{{.SHA}}")
BUILD_NAME = $(shell go run github.com/kellegous/glue/build/info@latest --format="{{.Name}}")

GO_MOD := $(shell go list -m)

ASSETS := \
	internal/ui/assets/index.html

BE_PROTOS := \
	reader.pb.go \
	reader_connect/reader.connect.go

FE_PROTOS := \
	ui/src/gen/reader_pb.ts

.PHONY: all clean develop nuke

.PRECIOUS: $(BE_PROTOS)

ALL: bin/reader

bin/%: cmd/%/main.go $(BE_PROTOS) $(ASSETS) $(shell find internal -name '*.go')
	go build -o $@ ./cmd/$*

bin/protoc:
	etc/download-protoc $(PROTOC_VERSION)

bin/protoc-gen-go:
	GOBIN="$(CURDIR)/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

bin/protoc-gen-connect-go:
	GOBIN="$(CURDIR)/bin" go install connectrpc.com/connect/cmd/protoc-gen-connect-go@$(PROTOC_GEN_CONNECT_GO_VERSION)

%.pb.go: %.proto bin/protoc-gen-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go_out=. \
		--go_opt=module=$(GO_MOD) \
		$<

reader_connect/reader.connect.go: reader.proto bin/protoc-gen-connect-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-connect-go=bin/protoc-gen-connect-go \
		--connect-go_out=. \
		--connect-go_opt=module=$(GO_MOD) \
		--connect-go_opt=package_suffix=_connect \
		$<

ui/src/gen/%_pb.ts: %.proto node_modules/.build bin/protoc
	mkdir -p $(dir $@)
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-es=node_modules/.bin/protoc-gen-es \
		--es_out=ui/src/gen \
		--es_opt=target=ts \
		$<

node_modules/.build:
	npm install
	touch $@

internal/ui/assets/index.html: node_modules/.build $(FE_PROTOS) $(shell find ui -type f)
	SHA="$(SHA)" BUILD_NAME="$(BUILD_NAME)" npm run build

develop: bin/reader
	bin/reader server --dev-mode=.:3020

clean:
	rm -rf bin internal/ui/assets

nuke: clean
	rm -rf node_modules
