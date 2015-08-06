#!/bin/bash

. ${COMMON}/common.sh

ssh_opts="-D $ssh_opts"

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

