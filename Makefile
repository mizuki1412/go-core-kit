BINARY=go-core-kit
VERSION=1.0.0
DATE=`date +%FT%T%z`
.PHONY: init

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