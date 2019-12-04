#!/bin/bash

ContainerNames=("5000/vulcand" "5000/thirdgate" "5000/adu" "5000/optimus" "5000/registry" "5000/nsq" "5000/etcd")
for name in ${ContainerNames[*]}
do
echo ${name}
result=`docker ps |grep ${name}`
# echo ${result}
if [ "${result}"x == "x" ]; then
continue
else
echo ${result:0:12}
sudo docker rm -f ${result:0:12}
fi
# echo ${result## }mm
done

