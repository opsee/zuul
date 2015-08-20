#!/bin/bash

usage() {
  cat <<EOF
Usage: $0

  -h            This message.
  -H HOSTNAME   Remote hostname or IP.
  -P PORT       Remote sshd port (default: 22).
  -p PORT       Local port to forward.
  -r PORT       Remote port to listen on (default: random).
  -u USERNAME   Username to connect with (default: zuul).
EOF

OPTIND=1

# connect -h some.host.name -p 22 -l 9022 -u username
while getopts "hH:P:p:r:u:" opt; do
  case "$opt" in
    h|\?)
      usage
      exit 0
      ;;
    H)
      host=$OPTARG
      ;;
    P)
      sshd_port=${OPTARG:-"22"}
      ;;
    p)
      local_port=$OPTARG
      ;;
    r)
      remote_port=${OPTARG:-"0"}
      ;;
    u)
      user=${OPTARG:-"zuul"}
  esac
done

if [ -z $SERVER_PUBLIC_KEY_OBJECT ] || [ -z $CLIENT_PRIVATE_KEY_OBJECT ]; then
  echo "Must specify SERVER_PUBLIC_KEY_OBJECT and CLIENT_PRIVATE_KEY_OBJECT"
  exit 1
fi

control_socket=/zuul/socket/zuul-${user}@${host}:${sshd_port}
# Technically, a multiplexer isn't required... But let's require it for
# now just to make sure we have it all working correctly.
if [ ! -s $control_socket ]; then
  echo "Control socket not found. Connecting directly..."
else
  echo "Found control socket at: ${control_socket}..."
  ssh_opts="-S $control_socket $ssh_opts"
fi

if [ -z $SERVER_PUBLIC_KEY_OBJECT ]; then
  echo "Must specify the S3 path to the public key..."
  exit 1
fi

get_object $SERVER_PUBLIC_KEY_OBJECT $server_public_key_path
if [ ! -r $server_public_key_path ]; then
  echo "Unable to read SSH public host key..."
  exit 1
fi
echo "$host $(< $server_public_key_path)" > $KNOWN_HOSTS_PATH


get_object $CLIENT_PRIVATE_KEY_OBJECT $client_private_key_path
if [ ! -r $client_private_key_path ]; then
  echo "Unable to read SSH private key..."
  exit 1
fi
chmod 0600 $client_private_key_path
