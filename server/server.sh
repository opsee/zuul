#!/bin/bash
# set -x

. ${COMMON}/common.sh

ssh_opts="-D $ssh_opts"
authorized_keys="${zuul_state}/authorized_keys"

sshd_config=/zuul/server/sshd_config
if [ ! -r $sshd_config ]; then
  echo "Cannot read sshd_config..."
  exit 1
else
  echo "Using config file ${sshd_config}..."
  ssh_opts="-f $sshd_config $ssh_opts"
fi

echo "Getting server private key..."
echo $server_private_key_path
get_object $SERVER_PRIVATE_KEY_OBJECT $server_private_key_path
if [ ! -r $server_private_key_path ]; then
  echo "Cannot read server private key ${server_private_key_path}..."
  exit 1
fi
chmod 600 $server_private_key_path
chown zuul:zuul $server_private_key_path
ssh_opts="-h $server_private_key_path $ssh_opts"

echo "Getting client public key..."
get_object $CLIENT_PUBLIC_KEY_OBJECT $authorized_keys
if [ ! -r $authorized_keys ]; then
  echo "Cannot read client public key ${authorized_keys}..."
  exit 1
fi
chmod 600 $authorized_keys
chown zuul:zuul $authorized_keys

/usr/sbin/sshd -d -f /zuul/server/sshd_config
