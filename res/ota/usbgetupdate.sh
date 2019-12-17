#! /bin/bash

autoupdate=1
url=$1
name=$2
if [ $autoupdate == 1 ]; then
    mkdir -p /ota
    cp $url $name.zip
    unzip $name.zip -d /ota/
    rm $name.zip
    if [ $? == 0 ];then
       echo "OK"
     else
       echo "ERROR"
       echo $?
    fi
else
   echo "no autoupdate"
fi
