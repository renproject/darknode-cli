MAIN_VERSION = $(shell cat ./VERSION | tr -d "[:space:]")
FULL_VERSION = ${MAIN_VERSION}

TARGET = ./darknode-cli-bin

# For information on flags: https://golang.org/cmd/link/
LDFLAGS = -s -w -X main.binaryVersion=${FULL_VERSION}

all: local

local: clean
	$(call build_local,./cmd)

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

