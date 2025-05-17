BINARY=go-core-kit
VERSION=2.1.4
DATE=`date +%FT%T%z`
.PHONY: build publish upgrade

default:
	@echo ${BINARY}
	@echo ${VERSION}
	@echo ${DATE}

upgrade:
	@go-mod-upgrade

publish:
	@git tag v${VERSION}
	@git push origin v${VERSION}

build:
	@GOOS=windows GOARCH=amd64 go build -tags timetzdata -o build/corekit.exe