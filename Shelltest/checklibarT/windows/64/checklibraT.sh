#!/bin/bash
IpListFile="./ip.txt"
ServerIp="192.168.100.235"
# 设备的启动脚本运行超限时间
LimitTime=300

./getDockerImageAndContainer.exe -serverIp=${ServerIp} -getIp=true -ipListFile=${IpListFile}

cat ${IpListFile} |while read Ip
do
    echo "================分割线================"
    echo ${Ip}
    # 测试IP能否ping通
    ip2=`ping ${Ip} -c1|head -2|tail -1|tr -s ' ' :|cut -d: -f4`
    if [ ${Ip}x == ${ip2}x ]; then
        # 获取设备reload_sensor.sh的运行情况
        ret1=`./sshrpc.exe -ip=${Ip} -p=false -l="ps -A -opid,stime,etime,args | grep reload_sensor.sh"`
        # 获取设备reload_sensor.sh运行时间
        reloadTime=`echo ${ret1} | grep -v grep |awk '{printf $3}'`
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
                    # ret2=`expr ${Array2} \* 60`
                    # ret=`expr ${ret1} + ${ret2} + ${Array3}`
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
        re4=`./sshrpc.exe -ip=${Ip} -p=false -l="ls /usr/bin/tegra_init.sh -l |cut -d \" \" -f 5"`
        re5=`./sshrpc.exe -ip=${Ip} -p=false -l="find /usr/bin/ -name tegra_init.sh"`
        # echo "/usr/bin/tegra_init.sh文件的大小为:${re4}"
        if [[ (-z ${re5} ) || (${re4} -lt 1754) ]];then
            echo "${Ip},设备的启动脚本损坏，请处理！"
        fi
        # 获取设备上正在运行的container和images以及设备版本
        ./getDockerImageAndContainer.exe -sensorIp=${Ip} -getDocker=true
    else
        echo "${Ip} ping 不通,请检查!"
    fi
done
