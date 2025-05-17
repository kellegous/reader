#!/bin/bash

set -euo pipefail

MINIFLUX_VERSION=2.2.8

install_miniflux() {
	local VERSION=$1
	local ARCH=$2
	curl -o /usr/local/bin/miniflux \
		-L \
		"https://github.com/miniflux/v2/releases/download/${VERSION}/miniflux-linux-${ARCH}"
	chmod +x /usr/local/bin/miniflux
}

export DEBIAN_FRONTEND=noninteractive

apt update

apt install -y postgresql sudo

apt clean

install_miniflux ${MINIFLUX_VERSION} $(go env GOARCH)

/bin/bash