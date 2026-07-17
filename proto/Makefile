SHELL:=/usr/bin/env bash

itick:
	GO_PKG="/pb"; \
	OPTS=$$( ls *.proto | \
		sed "s#^#--go_opt=M#" | \
		sed "s#\$$#=$$GO_PKG#"); \
	protoc --go_out=../server $$OPTS *.proto
