#!/bin/bash

while :
do
    count=`ps -a |grep main.go |wc -l`
    if [ ${count} != 2 ] ;then
        nohup go run main.go
    fi
    echo [ ${count} -eq 1 ]
    sleep 1s
done
