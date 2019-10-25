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
	@echo "package cri\n// Version of LXE\nconst Version = \"$(VERSION)\"" | gofmt > cri/version.go

.PHONY: build
build: version
	go build -v $(DEBUG) -o bin/$(EXECUTABLE) ./cmd/lxe

.PHONY: debug
debug: version
	go build -v -tags logdebug $(DEBUG) -o bin/$(EXECUTABLE) ./cmd/lxe

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
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run \
		--deadline 5m \
		--enable-all \
		--disable lll \
		--disable gochecknoglobals \
		--disable gochecknoinits \
		--disable godox \
		--disable funlen

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

.PHONY: package-clean
package-clean:
	rm -r package || true

.PHONY: package-deb-lxd-snap
package-deb-lxd-snap: build
	$(eval version:=$(shell make version))
	mkdir -p package/debian-lxd-snap/usr/bin
	
	objcopy --strip-debug --strip-unneeded --remove-section=.comment --remove-section=.note bin/$(EXECUTABLE) package/debian-lxd-snap/usr/bin/$(EXECUTABLE)
	cp -R fixtures/packaging/debian-lxd-snap/* package/debian-lxd-snap
	$(eval date:=$(shell date -R))
	VERSION="$(version)" DATE="$(date)" DOMAIN="$(DOMAIN)" envsubst < fixtures/packaging/debian-lxd-snap/usr/share/doc/$(DOMAIN)/changelog > package/debian-lxd-snap/usr/share/doc/$(DOMAIN)/changelog
	gzip -9 -S ".Debian.gz" package/debian-lxd-snap/usr/share/doc/$(DOMAIN)/changelog
	gzip -9 package/debian-lxd-snap/usr/share/man/man8/$(EXECUTABLE).8
	
	cd package/debian-lxd-snap; find . -type f -not -path './DEBIAN/*' -print | cut -c 3- | xargs md5sum > DEBIAN/md5sums
	VERSION="$(version)" INSTALLSIZE="\$$INSTALLSIZE" envsubst < fixtures/packaging/debian-lxd-snap/DEBIAN/control > package/debian-lxd-snap/DEBIAN/control
	du -s package/debian-lxd-snap | cut -f1 | awk '{print "Installed-Size: "$$1}' | sed -i -e '/$$INSTALLSIZE/{r /dev/stdin' -e 'd;}' package/debian-lxd-snap/DEBIAN/control

	chmod -R g-w package/debian-lxd-snap/
	fakeroot dpkg-deb -b package/debian-lxd-snap
	mv package/debian-lxd-snap.deb package/$(DOMAIN)_$(version).debian-lxd-snap.deb
	lintian --profile debian -i -I package/$(DOMAIN)_$(version).debian-lxd-snap.deb
