#! /bin/bash
name=$1
pid=`ps -ef | grep edgex-$name |grep -v color |grep -v grep | awk '{print $2}'`

if [ "${pid}" = "" ]
then
        echo "no java is alive"
else
        kill -9 ${pid}
fi

list_name="support-logging core-command core-data core-metadata export-client export-distro support-scheduler "   ###定义list
echo $name
if [[ "$list_name" =~ "$name" ]];then
    cd /app/edgex/$name
else
    cd /app/zap/$name
fi

LD_LIBRARY_PATH=/app/zeromq/lib exec -a edgex-$name ./$name >/app/log/edgex-$name.log 2>&1 &
