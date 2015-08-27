#!/bin/bash

. /zuul/common/common.sh

cat <<EOF >/dev/shm/up
${CUSTOMER_ID:-greg}
${VPN_PASSWORD:-crvo8u4B1Q1apn}
EOF
chown zuul:zuul /dev/shm/up
chmod 600 /dev/shm/up

openvpn --config /zuul/client/config
