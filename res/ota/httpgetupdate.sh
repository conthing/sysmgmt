#! /bin/bash

autoupdate=1
url=$1
name=$2
mkdir -p /ota
wget $url
unzip $name.zip -d /ota/
rm $name.zip
if [ $? == 0 ];then
   echo "OK"
else
   echo "ERROR"
   echo $?
fi

