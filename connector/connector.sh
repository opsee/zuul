#!/bin/bash

. ${COMMON}/common.sh
. ${COMMON}/client.sh

ssh_opts="-g -O forward -R ${local_port}:*:${remote_port} $ssh_opts"

ssh_config=/zuul/connector/ssh_config
if [ ! -r $ssh_config ]; then
  echo "Cannot read SSH config..."
  exit 1
else
  ssh_opts="-F $ssh_config $ssh_opts"
fi

echo "Forwarding..."
echo "LOCAL PORT: $local_port \
USERNAME: ${user} \
REMOTE HOST: ${host} \
REMOTE SSHD PORT: ${sshd_port}"

ssh $ssh_opts $host
