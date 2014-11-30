#!/bin/bash
if [[ -z "$1" ]]
then
/consul agent -server -bootstrap -data-dir /tmp/consul &
else
/consul agent -server -join "$1" -data-dir /tmp/consul &
fi
sleep 5
HOSTNAME=$(hostname)
consul-externalservice start --node "$HOSTNAME"
