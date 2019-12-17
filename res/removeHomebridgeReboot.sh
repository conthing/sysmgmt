#!/bin/bash

pidlist=`ps -ef | grep edgex-export-homebridge | grep -v "grep"|awk '{print $2}'`

if [ "$pidlist" = "" ]
then
echo "no pid alive"
else
echo "pidlist:$pidlist"

kill -9 $pidlist
sleep 1
fi

killall -9 homebridge

rm -rf /root/.homebridge/*

echo "now start homebright"

cd /app/zap/export-homebridge
LD_LIBRARY_PATH=/app/zeromq/lib exec -a edgex-export-homebridge ./export-homebridge >/app/log/edgex-export-homebridge.log 2>&1 &

sleep 2


/app/node/bin/homebridge -D >/app/log/edgex-homebridge 2>&1 &

echo "OK"
