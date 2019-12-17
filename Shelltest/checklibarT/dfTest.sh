#!/bin/bash
ret1=`./sshrpc.amd64 -ip=192.168.5.250 -p=false -l="df -h | awk 'NR==2{print}'"`
echo ${ret1}
# reloadTime=`echo ${ret1} | awk 'NR==2{print}' | awk '{print $4}'`
# echo ${reloadTime}
reloadTime2=`echo ${ret1} | awk '{print $5}'`
echo ${reloadTime2}
if [ ${reloadTime2}x == "100%x" ];then
echo "磁盘已满"
fi
