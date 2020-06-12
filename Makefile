# Simple Makefile for ercole agent

DESTDIR=build

all: linux windows

clean:
	rm -rf ercole-agent build ercole-agent.exe *.exe main

linux:
	GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent -tags linux -a

windows:
	GOOS=windows GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -o ercole-agent.exe -tags windows -a

nsis: windows
	makensis package/win/installer.nsi

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
