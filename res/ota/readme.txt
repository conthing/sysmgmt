准备工作
./prepare.sh

运行httpgetupdate.sh 
./httpgetupdate.sh 更新版本号 url链接地址 name解压文件名称
例：./httpgetupdate.sh 0.0.2 https://ota-insona.oss-cn-hangzhou.aliyuncs.com/app.zip app

#运行upgrade.sh 服务器会重启
#./upgrade.sh

运行recovery.sh
./recovery.sh 更新版本号 例：./recovery.sh 0.0.2
