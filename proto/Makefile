SHELL:=/usr/bin/env bash

PB_BIN=../client/node_modules/protobufjs-cli/bin
PB_JS=../client/src/pb

itick:
	GO_PKG="/pb"; \
	OPTS=$$(ls *.proto | \
		sed "s#^#--go_opt=M#" | \
		sed "s#\$$#=$$GO_PKG#"); \
	protoc --go_out=../server $$OPTS *.proto

client:
	$(PB_BIN)/pbjs -t static-module -w es6 --es6 --force-number -o ${PB_JS}/pb.js vapid.proto
	$(PB_BIN)/pbts -o ${PB_JS}/pb.d.ts ${PB_JS}/pb.js
