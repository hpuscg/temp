#!/bin/bash

IP=$(hostname -I | cut -f1 -d' ')

## modify eventserver service
sudo docker tag 192.168.5.46:5000/armhf-eventserver:1.7.7 $IP:5000/armhf-eventserver:1.7.7
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-eventserver/cur_version -XPUT -d value=1.7.7
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-eventserver/to_upgrade -XPUT -d value='true'
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-eventserver/1.7.7/image -XPUT -d value='192.168.5.46:5000/armhf-eventserver:1.7.7"
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-eventserver/1.7.7/script -XPUT -d value="sudo docker run -tid --name eventserver --memory=800m --restart=always --net=host -v /tmp:/tmp -v /data/tf/eventserver:/data/slice -v /data/eventserver:/data/event -v /etc/localtime:/etc/localtime:ro IMAGE ./eventserver.arm -etcdserver=http://127.0.0.1:4001 -mode=fat -report_interval=600 -log_dir_deepglint /tmp/ -mem=750"

## modify pioneer service
sudo docker tag 192.168.5.46:5000/armhf-pioneer:1.5.7 $IP:5000/armhf-pioneer:1.5.7
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-pioneer/cur_version -XPUT -d value=1.5.7
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-pioneer/to_upgrade -XPUT -d value='true'
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-pioneer/1.5.7/image -XPUT -d value='192.168.5.46:5000/armhf-pioneer:1.5.7"
curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-pioneer/1.5.7/script -XPUT -d value="sudo docker run -i -t --rm --privileged --net=host IMAGE /bin/bash init.sh"

## modify system version 版本号
curl -L http://127.0.0.1:2379/v2/keys/config/global/cur_version -XPUT -d value="V2.13.190101R"
