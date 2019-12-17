#!/bin/bash

killall -9 homebridge

echo "now start homebright"

/app/node/bin/homebridge -D >/app/log/edgex-homebridge 2>&1 &

echo "OK"
