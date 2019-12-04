#!/bin/bash

# SERVER=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/server_addr`
SERVER="192.168.4.42"

UPGRADE_SERVICE=`curl --request GET --url http://$SERVER:8008/api/config?key=/config/global/upgrade_service`


for TEMP in $UPGRADE_SERVICE
    do
        echo $TEMP
    done
