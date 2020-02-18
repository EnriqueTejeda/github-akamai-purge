# Helpers
ROOT_DIR:=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build
build:
	@docker build . -t etejeda/github-akamai-purge:latest
