#!/bin/bash

set -euo pipefail

MINIFLUX_VERSION=2.2.13

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

apt update

apt install -y curl gnupg

curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor -o /usr/share/keyrings/postgresql-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/postgresql-keyring.gpg] http://apt.postgresql.org/pub/repos/apt/ bookworm-pgdg main" | tee /etc/apt/sources.list.d/postgresql.list

apt update

apt install -y postgresql-14 sudo ca-certificates

apt clean

install_miniflux ${MINIFLUX_VERSION} ${ARCH}