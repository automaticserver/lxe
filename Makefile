DOMAIN=lxe
EXECUTABLE=lxe

LXDSOCKETFILE ?= /var/lib/lxd/unix.socket
LXESOCKETFILE ?= /var/run/lxe.sock
LXDSOCKET=unix://$(LXDSOCKETFILE)
LXESOCKET=unix://$(LXESOCKETFILE)
LXELOGFILE ?= /var/log/lxe.log

VERSION=$(shell (git describe --tags --dirty --always --exact-match --match 'v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null || (echo -n "v0.0.0-"; git describe --dirty --always)) | cut -c 2- )
PACKAGENAME=$(shell echo "$${PWD\#"$$GOPATH/src/"}")

GO111MODULE=on

.PHONY: all
all: build test lint

.PHONY: version
version:
	@echo "$(VERSION)"
	@echo -e "package cri\n// Version of LXE\nconst Version = \"$(VERSION)\"" | gofmt > cri/version.go

.PHONY: build
build: version
	go build -o bin/$(EXECUTABLE) -gcflags=all='-N -l' ./cmd/lxe

.PHONY: generate
generate:
	go mod tidy
	go generate ./...

# $(GOPATH)/bin/critest:
# # 	go get -v -u "github.com/kubernetes-incubator/cri-tools/cmd/critest"
# 	make -C $(GOPATH)/src/github.com/kubernetes-incubator/cri-tools critest

.PHONY: test
test:
	go test ./... -coverprofile go.coverprofile
	go tool cover -func go.coverprofile
# | sed 's|^$(shell go list -m)/||'

.PHONY: lint
lint:
	go mod download
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: integration-test
integration-test: critest

.PHONY: clean
clean: package-clean
	rm -r bin || true

.PHONY: checklxd
checklxd:
	@test -e $(LXDSOCKETFILE) || (echo "Socket $(LXDSOCKETFILE) not found! Is LXD running?" && false)
	@test -r $(LXDSOCKETFILE) || (echo "Socket $(LXDSOCKETFILE) not accessible! Can this user read it?" && false)

.PHONY: prepareintegration
prepareintegration:
	lxc image copy images:alpine/edge local: --alias busybox \
		--alias gcr.io/cri-tools/test-image-latest:latest \
		--alias gcr.io/cri-tools/test-image-digest@sha256:9179135b4b4cc5a8721e09379244807553c318d92fa3111a65133241551ca343

.PHONY: critest
critest: checklxd $(GOPATH)/bin/critest
	$(GOPATH)/bin/critest -runtime-endpoint	$(LXESOCKET) -image-endpoint $(LXESOCKET)

.PHONY: cribench
cribench: checklxd default prepareintegration $(GOPATH)/bin/critest
	(./bin/$(EXECUTABLE) --socket $(LXESOCKETFILE) --lxd-socket $(LXDSOCKETFILE) --logfile $(LXELOGFILE) &)
	$(GOPATH)/bin/critest -benchmark -runtime-endpoint $(LXESOCKET) -image-endpoint $(LXESOCKET)

run: checklxd build
	./bin/$(EXECUTABLE) --debug --socket $(LXESOCKETFILE) --lxd-socket $(LXDSOCKETFILE) --logfile $(LXELOGFILE)
