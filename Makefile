BINARY=go-core-kit
VERSION=1.2.3
DATE=`date +%FT%T%z`
.PHONY: init build

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}

init:
	@go generate
	@echo "[ok] generate"

upgrade:
	@go-mod-upgrade

publish:
	@git tag v${VERSION}
	@git push origin v${VERSION}

build:
	@GOOS=windows GOARCH=amd64 go build -o build/corekit.exe