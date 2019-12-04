#!/bin/bash

# EVENTSERVER=`sudo docker ps -a | awk '{ print $NF }' | grep eventserver`
# if [ "$EVENTSERVER"x = "eventserver"x ]; then
#	#statements
#	sudo docker stop eventserver
#	sleep 1s
#	sudo docker start eventserver
# else
#	KEY="config/image/armhf-eventserver"

#	IN_USE=`etcdctl --peers 127.0.0.1:4001 get $KEY/in_use`

#	if [ "$IN_USE"x != "true"x ]; then
#		#statements
#		logger "Eventserver false."
#		exit
#	fi

#	CUR_VERSION=`etcdctl --peers 127.0.0.1:4001 get $KEY/cur_version`
#	CUR_IMAGE=`etcdctl --peers 127.0.0.1:4001 get $KEY/$CUR_VERSION/image`
#	CUR_SCRIPT=`etcdctl --peers 127.0.0.1:4001 get $KEY/$CUR_VERSION/script`

#	sudo docker pull $CUR_IMAGE >/dev/null 2>&1

#	${CUR_SCRIPT/IMAGE/$CUR_IMAGE}
# fi

VODSERVER=`sudo docker ps -a | awk '{ print $NF }' | grep vodserver`
if [ "$VODSERVER"x = "vodserver"x ]; then
	#statements
	sudo docker stop vodserver 
	sleep 1s
	sudo docker start vodserver
else
	KEY="config/image/armhf-vodserver"

	IN_USE=`etcdctl --peers 127.0.0.1:4001 get $KEY/in_use`

	if [ "$IN_USE"x != "true"x ]; then
		#statements
		logger "Vodserver false."
		exit
	fi

	CUR_VERSION=`etcdctl --peers 127.0.0.1:4001 get $KEY/cur_version`
	CUR_IMAGE=`etcdctl --peers 127.0.0.1:4001 get $KEY/$CUR_VERSION/image`
	CUR_SCRIPT=`etcdctl --peers 127.0.0.1:4001 get $KEY/$CUR_VERSION/script`

	sudo docker pull $CUR_IMAGE >/dev/null 2>&1

	${CUR_SCRIPT/IMAGE/$CUR_IMAGE}
fi
