#!/bin/bash
echo "OK"
ntptype=$1
date=$2
url=$3
echo $ntptype
echo $date
echo $url
if [ $ntptype = "ntp" ];then
	echo 1
    timedatectl set-ntp true
    ntpdate $url
    hwclock -w
else
	echo 2
    timedatectl set-ntp false
    date -s "$date"
    hwclock -w
fi
