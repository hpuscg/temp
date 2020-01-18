#!/bin/bash
IpListFile="./ip.txt"
ServerIp="192.168.100.235"
# 设备的启动脚本运行超限时间
LimitTime=300

# ./getDockerImageAndContainer.amd64 -serverIp=${ServerIp} -getIp=true -ipListFile=${IpListFile}
i=0
cat ${IpListFile} |while read Ip
do
    ((i++))
    echo "================${i}================"
    echo ${Ip}
    # 测试IP能否ping通
    ip2=`ping ${Ip} -c1|head -2|tail -1|tr -s ' ' :|cut -d: -f4`
    if [ ${Ip}x != ${ip2}x ]; then
        echo "${Ip} ping 不通,请检查!"
        continue
    fi
    # 检查设备是否配置了网管服务器
    server_addr=`./sshrpc.amd64 -ip=${Ip} -p=false -l="etcdctl get /config/global/server_addr"`
    if [ ${server_addr}x == "x" ]; then
        echo "${Ip} 未配置网管服务器"
    fi
    # 检查设备的磁盘空间
    dfret=`./sshrpc.amd64 -ip=${Ip} -p=false -l="df | awk 'NR==2{print}'"`
    dfResult=`echo ${dfret} | awk '{print $5}'`
    if [ ${dfResult}x = "100%x" ];then
        echo "磁盘已满"
    else
        echo "磁盘已使用：${dfResult}"
    fi
    # 获取设备reload_sensor.sh的运行情况
    rets=`./sshrpc.amd64 -ip=${Ip} -p=false -l="ps -A -opid,stime,etime,args | grep reload_sensor.sh"`
    # 获取设备reload_sensor.sh运行时间
    reloadTime=`echo ${rets} | grep -v grep |awk '{printf $3}'`
    if [ ${reloadTime}x != "x" ];then
        Array1=`echo ${reloadTime} |cut -f 1 -d ":"`
        Array2=`echo ${reloadTime} |cut -f 2 -d ":"`
        Array3=`echo ${reloadTime} |cut -f 3 -d ":"`
        # 判断启动脚本运行时间是否超过1小时
        if [ ${Array3}x == "x" ];then
            ret1=`expr ${Array1} \* 60`
            ret=`expr ${ret1} + ${Array2}`
        else
            Array11=`echo ${Array1} |cut -f 1 -d "-"`
            Array12=`echo ${Array1} |cut -f 2 -d "-"`
            # 判断启动脚本运行时间是否超过一天
            if [ ${Array12}x == "x" ];then
                ret=`expr ${Array11} \* 3600 + ${Array2} \* 60 + ${Array3}`
            # 启动脚本运行时间超过一天
            else
                ret0=`expr ${Array11} \* 86400`
                ret1=`expr ${Array12} \* 3600`
                ret2=`expr ${Array2} \* 60`
                ret=`expr ${ret0} + ${ret1} + ${ret2} + ${Array3}`
            fi
        fi
        # 判断启动脚本运行时间是否超限
        if [ ${ret} -ge ${LimitTime} ];then
            echo "reload_sensor.sh 卡死"
        fi
    fi
    # 测试设备启动脚本是否损坏
    re4=`./sshrpc.amd64 -ip=${Ip} -p=false -l="ls /usr/bin/tegra_init.sh -l |cut -d \" \" -f 5"`
    re5=`./sshrpc.amd64 -ip=${Ip} -p=false -l="find /usr/bin/ -name tegra_init.sh"`
    if [[ (-z ${re5} ) || (${re4} -lt 1754) ]];then
        echo "${Ip},设备的启动脚本损坏，请处理！"
    fi
    # 判断设备上的服务是否正常
    ./getDockerImageAndContainer.amd64 -sensorIp=${Ip} -serverTest=true
    # 获取设备上正在运行的container和images以及设备版本
    ./getDockerImageAndContainer.amd64 -sensorIp=${Ip} -getDocker=true
done
