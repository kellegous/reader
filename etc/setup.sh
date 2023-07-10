#!/bin/bash

set -euo pipefail

MINIFLUX_VERSION=2.0.45

install_miniflux() {
	local VERSION=$1
	local ARCH
	case "$2" in
		"x86_64")
		ARCH="amd64"
		;;
		"aarch64")
		ARCH="arm64"
		;;
		*)
	esac
	curl -o /usr/local/bin/miniflux \
		-L \
		"https://github.com/miniflux/v2/releases/download/${VERSION}/miniflux-linux-${ARCH}"
	chmod +x /usr/local/bin/miniflux
}

export DEBIAN_FRONTEND=noninteractive

# Check for supported architectures
ARCH=$(uname -m)
case "$ARCH" in
	"x86_64"|"aarch64")
	;;
	*)
	echo "Unsupported machine arch $(uname -m)"
	exit 1
	;;
esac

apt-get update

apt-get install -y postgresql postgresql-contrib sudo ca-certificates curl

apt-get clean

install_miniflux ${MINIFLUX_VERSION} ${ARCH}