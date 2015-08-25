#!/bin/bash

. /zuul/common/common.sh

get_encrypted_object dev/opsee-key.pem /zuul/state/server.key
get_encrypted_object dev/opsee.crt /zuul/state/server.crt

/usr/sbin/openvpn --config /zuul/server/server.conf
