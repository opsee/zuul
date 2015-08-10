#!/bin/bash

. ${COMMON}/common.sh

ssh_opts="-D $ssh_opts"
authorized_keys="${STATE}/authorized_keys"

sshd_config=/zuul/server/sshd_config
if [ ! -r $sshd_config ]; then
  echo "Cannot read sshd_config..."
  exit 1
else
  echo "Using config file ${sshd_config}..."
  ssh_opts="-f $sshd_config $ssh_opts"
fi

echo "Getting server private key..."
get_object $SERVER_PRIVATE_KEY_OBJECT $SERVER_PRIVATE_KEY_PATH
if [ ! -r $SERVER_PRIVATE_KEY_PATH ]; then
  echo "Cannot read server private key ${SERVER_PRIVATE_KEY_PATH}..."
  exit 1
fi
chmod 600 $SERVER_PRIVATE_KEY_PATH
ssh_opts="-h $SERVER_PRIVATE_KEY_PATH $ssh_opts"

echo "Getting client public key..."
get_object $CLIENT_PUBLIC_KEY_OBJECT $authorized_keys
if [ ! -r $CLIENT_PUBLIC_KEY_PATH ]; then
  echo "Cannot read client public key ${authorized_keys}..."
  exit 1
fi
chmod 600 $authorized_keys
mv 

sshd -f /zuul/server/sshd_config
