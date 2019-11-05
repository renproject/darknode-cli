MAIN_VERSION = $(shell cat ./VERSION | tr -d "[:space:]")
BRANCH = $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT_HASH = $(shell git describe --always --long)
FULL_VERSION = ${MAIN_VERSION}-${BRANCH}-${COMMIT_HASH}

TARGET = ./darknode_bin

# For information on flags: https://golang.org/cmd/link/
LDFLAGS = -s -w -X main.binaryVersion=${FULL_VERSION}

all: local

local:
	$(call build_local,./cmd/darknode)

version:
	@ echo ${FULL_VERSION}

target_name:
	@ echo "${TARGET}"

clean:
	rm -rf "${TARGET}"

define build_local
	go build -o ${TARGET} -ldflags="${LDFLAGS}" $(1)
endef

.PHONY: all local version clean target_name

