BINARY=go-core-kit
VERSION=2.1.6
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
	@GOOS=windows GOARCH=amd64 go build -tags timetzdata -trimpath -ldflags="-s -w" -o build/corekit.exe

# go build 参数说明
## -tags timetzdata 携带时区
## -trimpath 移除所有记录在可执行文件中的绝对文件路径
## -ldflags="-s -w" 移除调试符号信息, 会影响dlv调试
