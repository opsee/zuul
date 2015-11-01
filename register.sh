#!/usr/bin/env bash
set -a

if [ -x /opt/bin/ec2-env ]; then
  /opt/bin/ec2-env > /gozer/state/ec2.environment
  if [ $? -eq 0 ]; then
    eval "$(< /gozer/state/ec2.environment)"
  fi
fi

/zuul/bin/register
