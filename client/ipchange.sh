#!/bin/bash
set -a
set -e

ip=$1

echo $ip > /zuul/state/ip
