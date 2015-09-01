#!/bin/bash
set -a

if [ -x /opt/bin/ec2-env ]; then
  /opt/bin/ec2-env > /zuul/state/ec2.environment
  if [ $? -eq 0 ]; then
    eval "$(< /zuul/state/ec2.environment)"
  fi
fi

NSQD_HOST="nsqd.opsy.co:4150"

/zuul/gozer/bin/register
