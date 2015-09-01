#!/bin/bash
set -e
set -a

if [ -x /opt/bin/ec2-env ]; then
  eval "$(/opt/bin/ec2-env)"
fi

NSQD_HOST="nsqd.opsy.co:4150"

/zuul/gozer/bin/register
