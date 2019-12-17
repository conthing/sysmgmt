#!/bin/bash

ip_file=/etc/init.d/ipaddr.sh
#zt=`grep "static" $ip_file|wc -l`
netname=$1
nettype=$2
address=$3
netmask=$4
gateway=$5


sed -i '/ifconfig/d' $ip_file
sed -i '/udhcp/d' $ip_file
sed -i '/route/d' $ip_file
if [ $nettype == "dhcp"  ];then
cat >>$ip_file<< EOF
udhcpc -i $netname -R
EOF

else

cat >>$ip_file<< EOF
ifconfig $netname down
ifconfig $netname $address netmask $netmask up
route add default gw $gateway
EOF

fi

sleep 2
if [ $nettype == "dhcp"  ];then
	udhcpc -i $netname -R
	echo "dhcp"
else
	kill -9 `ps -ef | grep udhcpc | awk '{print $2}'`
	ifconfig $netname down
	ifconfig $netname $address netmask $netmask up
	route add default gw $gateway
	echo "static $address"
fi
echo "OK"
