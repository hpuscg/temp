#!/bin/bash
ips=("")
for Ip in ${ips[@]}
do
    echo ${Ip}
    ./setIotInfo.linux -ip=${Ip}
done