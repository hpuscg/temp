#!/bin/bash


ips=("" "" "")

for Ip in ${ips[@]}
do
    echo ${Ip}
    ./sshrpc.exe -ip ${Ip} -l "ls -lh /usr/bin/tegra_init.sh"
    ./sshrpc.exe -ip ${Ip} -l "ls -lh /data/shell/_usrbin"
    echo "===================分割线======================"
done
