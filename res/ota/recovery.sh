#! /bin/bash

if [ ! -d /ota/app ];then
    echo "app下载解压失败"
    exit
fi

update-rc.d -f recovery.sh remove
rm /etc/init.d/recovery.sh

for((i=0;i<12;i++));
do
sleep 5m
result=$(curl http://localhost:48081/api/v1/ping)
if [[ $result == "pong" ]];then
        echo "OK"
else
	echo "not OK"
	if [ $i > 3 ];then
	    cd /app
		./shutdown.sh
		cd /sysmgmt/ota/
		./callback.sh
		cd /app
		./startup.sh
		break
	fi
fi
done