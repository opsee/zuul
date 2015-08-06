#!/bin/bash

echo "Setting up runtime..."

PATH=$PATH:/opt/bin

check_env_var() {
  local var=$1
  if [ -z "${!var}" ]; then
    echo "Must supply $var environment variable"
    exit 1
  fi
}

get_object() {
  local key=$KEY_ALIAS
  local bucket=$KEY_BUCKET
  local obj=$1
  local target=$2
  
  if [ -z "$obj" ] || [ -z "$target" ]; then
    echo "get_object requires two arguments"
    exit 1
  fi

  s3kms -k $key get -b $bucket -o $obj > $target
}

required="KEY_ALIAS \
KEY_BUCKET \
SERVER_PUBLIC_KEY_OBJECT \
AWS_DEFAULT_REGION \
AWS_ACCESS_KEY_ID \
AWS_SECRET_ACCESS_KEY"

STATE=${STATE:-"/zuul/state"}
KNOWN_HOSTS_PATH=${STATE}/ssh_known_hosts
SERVER_PRIVATE_KEY_PATH=${STATE}/ssh_server_key
SERVER_PUBLIC_KEY_PATH=${SERVER_PRIVATE_KEY_PATH}.pub
CLIENT_PRIVATE_KEY_PATH=${STATE}/ssh_client_key
CLIENT_PUBLIC_KEY_PATH=${CLIENT_PRIVATE_KEY_PATH}.pub

for v in $required; do
  check_env_var $v
done

if [ ! -d $state_dir ]; then
  mkdir -p $state_dir
fi

ssh_opts="-v"
