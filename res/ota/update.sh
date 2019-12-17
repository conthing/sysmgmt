#!/bin/bash

type=$1
filepath=$2
filedir=$3
echo $type
echo $filepath
echo $filedir

if [ ! -d /ota/app ];then
    echo "app下载解压失败"
    exit
fi
if [ $type = "ADD" ];then
    if [ ! -d /app/$filedir ];then
        mkdir -p /app/$filedir
    fi
    cd /ota/app
    mv $filepath /app/$filepath
elif [ $type = "MODIFIED" ];then
    echo "MODIFIED"
    if [ ! -d /app.bak/$filedir ];then
        mkdir -p /app.bak/$filedir
    fi
    cd /ota/app
    mv /app/$filepath /app.bak/$filepath
    mv $filepath /app/$filepath
elif [ $type = "DELETE" ];then
    echo "DELETE"
    if [ ! -d /app.bak/$filedir ];then
        mkdir -p /app.bak/$filedir
    fi
    cd /ota/app
    rm -rf /app.bak/$filepath
    mv /app/$filepath /app.bak/$filepath
fi