#!/bin/bash

. /zuul/common/common.sh

cat <<EOF >/dev/shm/up
${BASTION_ID}
${VPN_PASSWORD}
EOF
chown zuul:zuul /dev/shm/up
chmod 600 /dev/shm/up

openvpn --config /zuul/client/config
