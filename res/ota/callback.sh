#! /bin/bash

type=$1
filepath=$2
filedir=$3

if [ $type = "ADD" ];then
    rm -rf /app/$file
    if [ ! -d /app/$filedir ];then
        mkdir -p /app/$filedir
    fi
    cd /ota/app
    mv $filepath /app/$filepath
elif [ $type = "MODIFIED" ];then
    rm -rf /app/$file
    if [ ! -d /app.bak/$filedir ];then
       mkdir -p /app.bak/$filedir
    fi
    mv /app.bak/$filepath /app/$filepath
elif [ $type = "DELETE" ];then
    if [ ! -d /app.bak/$filedir ];then
        mkdir -p /app.bak/$filedir
    fi
    mv /app.bak/$filepath /app/$filepath
fi


