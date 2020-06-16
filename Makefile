# Simple Makefile for ercole agent

DESTDIR=build

all: ercole-agent

default: ercole-agent

clean:
	rm -rf ercole-agent build ercole-agent.exe *.exe

ercole-agent:
	GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent -a

windows:
	GOOS=windows GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent.exe -a

nsis: windows
	makensis package/win/installer.nsi

rhel5:
	docker run --rm -it -v "$$PWD":/go/src/github.com/ercole-io/ercole-agent -w /go/src/github.com/ercole-io/ercole-agent golang:1.3 go build -tags rhel5

test:
	go test ./...
	go test -tags rhel5 ./...

install: all install-fetchers install-bin install-bin install-config install-scripts

install-fetchers:
	install -d $(DESTDIR)/fetch
	cp -rp fetch/* $(DESTDIR)/fetch
	rm $(DESTDIR)/fetch/win.ps1

install-bin:
	install -m 755 ercole-agent $(DESTDIR)/ercole-agent
	install -m 755 package/ercole-setup $(DESTDIR)/ercole-setup

install-scripts:
	install -d $(DESTDIR)/sql
	install -m 644 sql/*.sql $(DESTDIR)/sql

install-config:
	install -m 644 config.json $(DESTDIR)/config.json
