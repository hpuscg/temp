#!/bin/bash

if [ $# -gt 0 ]; then
	IP=$1
else
    echo ./upinit.sh sersor_ip
    exit
fi

./sshrpc -ip $IP -c 'etcdctl set /config/libra/data/enable_bodycutboard "1"'
./sshrpc -ip $IP -c 'etcdctl set /config/libra/data/cutboard_interval "20000"'
./sshrpc -ip $IP -c 'etcdctl set /config/libra/data/image_buffer_size "100"'
./sshrpc -ip $IP -c 'etcdctl set /config/libra/data/enable_color_tracking "1"'
./sshrpc -ip $IP -c 'etcdctl set /config/libra/data/enable_framecutboard "1"'
./sshrpc -ip $IP -c 'reload_sensor.sh'

exit

