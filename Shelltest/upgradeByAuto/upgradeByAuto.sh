#!/bin/bash
IpListFile="./ip.txt"
ServerIp="192.168.100.235"

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
    # 检查设备的磁盘空间
    dfret=`./sshrpc.amd64 -ip=${Ip} -p=false -l="df | awk 'NR==2{print}'"`
    dfResult=`echo ${dfret} | awk '{print $5}'`
    if [[ ${dfResult%%%} > 80 ]];then
        echo "${Ip} 磁盘已使用 ${dfResult}"
        continue
    fi
    # 测试设备启动脚本是否损坏
    re4=`./sshrpc.amd64 -ip=${Ip} -p=false -l="ls /usr/bin/tegra_init.sh -l |cut -d \" \" -f 5"`
    re5=`./sshrpc.amd64 -ip=${Ip} -p=false -l="find /usr/bin/ -name tegra_init.sh"`
    if [[ (-z ${re5} ) || (${re4} -lt 1754) ]];then
        echo "${Ip},设备的启动脚本损坏，请处理！"
        # 将设备启动脚本移动到/usr/bin下
        ./sshrpc.amd64 -ip=${Ip} -p=false -l="cp /data/shell/_usrbin/tegra_init.sh /usr/bin/"
    fi
    # 检查设备是否配置了网管服务器
    server_addr=`./sshrpc.amd64 -ip=${Ip} -p=false -l="etcdctl get /config/global/server_addr"`
    if [ ${server_addr}x == "x" ]; then
        # 未配置网管服务器，较时，配置网管服务器
        ./upgradeByAuto.amd64 -sensorIp=${Ip}
    else
        # TODO 检查设备版本，是否需要预升级
        package_size=`./sshrpc.amd64 -ip=${Ip} -p=false -l="ls /data/shell/service/hookshell/switch_package.sh -l |cut -d \" \" -f 5"`
        ./upgradeByAuto.amd64 -sensorIp=${Ip} -hasServerAddr=true
    fi
done