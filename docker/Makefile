SHELL:=/usr/bin/env bash

PARENT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)

build:
	sudo BUILDX_GIT_INFO=0 docker build \
		--provenance=false \
		--progress=plain \
		--tag stock-image \
		--build-arg "GIT_COMMIT=$$(../build/git-hash.sh)" \
		-f Dockerfile \
		..

run:
	sudo docker run \
		--mount type=bind,source=$(PARENT_DIR)/static,target=/static \
		-p 127.0.0.1:22721:80 \
		--name stock-debug --rm \
		stock-image

tanya:
	./deploy-prod.sh
