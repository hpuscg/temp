#!/bin/bash
ips=("")
Value="1"
for Ip in ${ips[@]}
do
    echo ${Ip}
    ./startCutbroad.linux -ip=${Ip} -value=${Value}
done