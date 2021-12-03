# Simple Makefile for ercole agent

DESTDIR=build

.PHONY: all
all: ercole-agent

.PHONY: default
default: ercole-agent

.PHONY: clean
clean:
	rm -rf ercole-agent build ercole-agent.exe *.exe
	find . -name "fake_*_test.go" -exec rm "{}" \;
	go generate ./...
	go clean -testcache

ercole-agent:
	GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent -a

windows:
	GOOS=windows GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent.exe -a

nsis: windows
	makensis package/win/installer.nsi

.PHONY: install
install: all install-fetchers install-bin install-bin install-config install-scripts

.PHONY: install-fetchers
install-fetchers:
	install -d $(DESTDIR)/fetch
	cp -rp fetch/linux $(DESTDIR)/fetch

.PHONY: install-bin
install-bin:
	install -m 755 ercole-agent $(DESTDIR)/ercole-agent
	install -m 755 package/ercole-setup $(DESTDIR)/ercole-setup

.PHONY: install-scripts
install-scripts:
	install -d $(DESTDIR)/sql
	cp -rp sql/* $(DESTDIR)/sql

.PHONY: install-config
install-config:
	install -m 644 config.json $(DESTDIR)/config.json

.PHONY: test
test:
	go test ./...
