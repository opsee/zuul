#!/bin/bash

set -a

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
  local bucket=$KEY_BUCKET
  local obj=$1
  local target=$2

  echo "Copying ${obj} from KMS to ${target} locally."

  if [ -z "$obj" ] || [ -z "$target" ]; then
    echo "get_object requires two arguments"
    exit 1
  fi

  s3kms get -b $bucket -o $obj > $target
}

required="KEY_ALIAS \
KEY_BUCKET \
AWS_DEFAULT_REGION \
AWS_ACCESS_KEY_ID \
AWS_SECRET_ACCESS_KEY"

known_hosts_path=${zuul_state}/ssh_known_hosts
zuul_state=${zuul_state:-"/zuul/state"}
server_private_key_path=${zuul_state}/ssh_server_key
server_public_key_path=${server_private_key_path}.pub
client_private_key_path=${zuul_state}/ssh_client_key
client_public_key_path=${client_private_key_path}.pub

for v in $required; do
  check_env_var $v
done

if [ ! -d $state_dir ]; then
  mkdir -p $state_dir
fi

ssh_opts="-v"
