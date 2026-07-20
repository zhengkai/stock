SHELL:=/usr/bin/env bash

-include ./build/config.ini

start:
	./build/run-server.sh $(type)

stop:
	./build/stop-server.sh $(type)
