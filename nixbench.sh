#!/bin/sh

set -e

case "$(uname -s)" in
    Linux)      os=linux;;
    Darwin)     os=darwin;;
    FreeBSD)    os=freebsd;;
    SunOS)      os=solaris;;
    DragonFly)  os=dragonfly;;
    NetBSD)     os=netbsd;;
    OpenBSD)    os=openbsd;;
esac

if [[ -z "$os" ]]; then
  echo "Unsupported operating system: $(uname -s)"
  exit 1
fi

case "$(uname -m)" in
    amd64)      arch=amd64;;
    x86_64)     arch=amd64;;
    i686)       arch=386;;
    i386)       arch=386;;
    armv6l)     arch=arm;;
    armv7l)     arch=arm;;
esac

if [[ -z "$arch" ]]; then
  echo "Unsupported architecture: $(uname -m)"
  exit 1
fi

url=$(curl -s https://api.github.com/repos/jgillich/nixbench/releases/latest | grep browser_download_url | cut -d '"' -f 4 | grep $os-$arch | xargs)

curl -LsSo /tmp/nixbench $url

chmod +x /tmp/nixbench

/tmp/nixbench
