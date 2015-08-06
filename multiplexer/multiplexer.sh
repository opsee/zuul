#!/bin/bash

. ${COMMON}/common.sh
. ${COMMON}/client.sh
# provides: host sshd_port local_port remote_port user

ssh_config="/zuul/multiplexer/ssh_config"
if [ ! -r $ssh_config ]; then
  echo "Cannot read SSH config..."
  exit 1
else
  ssh_opts="-F $ssh_config $ssh_opts"
fi

required="host port user"
for v in $required; do
  check_env_var $v
done

ssh_opts="-M -T -p $port -l $user"

ssh $ssh_opts $host

