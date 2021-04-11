#!/bin/bash
echo "ModifyTime $1 $2 $3"
ntptype=$1
date=$2
url=$3
if [ $ntptype = "ntp" ];then
	echo "timedatectl set-ntp true"
    timedatectl set-ntp true
	echo "ntpdate $url"
    ntpdate $url
	echo "hwclock -w"
    hwclock -w
else
	echo "timedatectl set-ntp false"
    timedatectl set-ntp false
	echo "date -s $date"
    date -s "$date"
	echo "hwclock -w"
    hwclock -w
fi
