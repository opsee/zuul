#!/bin/bash

. /zuul/common/common.sh

get_encrypted_object dev/opsee-key.pem /zuul/state/server.key
get_encrypted_object dev/opsee.crt /zuul/state/server.crt

if [ -z "$1" ]; then
  echo "Must specify a /24 network to start a server."
  exit 1
fi

/usr/sbin/openvpn --config /zuul/server/server.conf --server $1 255.255.255.0
