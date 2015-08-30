#!/bin/bash
# set -x
set -a

echo "Setting up runtime..."

if [ ! -d /zuul/state ]; then
  mkdir -p /zuul/state
fi

/opt/bin/ec2-env > /zuul/state/environment
if [ $? -eq 0 ]; then
  eval "$(< /zuul/state/environment)"
fi

get_encrypted_object() {
  local bucket=opsee-keys
  local obj=$1
  local target=$2

  echo "Copying ${obj} from KMS to ${target} locally."

  if [ -z "$obj" ] || [ -z "$target" ]; then
    echo "get_object requires two arguments"
    exit 1
  fi

  s3kms -r "us-west-1" get -b $bucket -o $obj > $target
  chmod 600 $target
  if [ ! -r $target ]; then
    echo "Unable to read $target"
    exit 1
  fi
}

get_encrypted_object dev/ca.crt /zuul/state/ca.crt
get_encrypted_object dev/tls-auth.key /zuul/state/tls-auth.key
