#!/bin/bash

. /zuul/common/common.sh

cat <<EOF >/dev/shm/up
${CUSTOMER_ID:-cliff}
${VPN_PASSWORD:-cliff}
EOF
chown zuul:zuul /dev/shm/up
chmod 600 /dev/shm/up

openvpn --config /zuul/client/client.conf
