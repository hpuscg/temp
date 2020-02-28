#!/bin/bash

if [ -f /run/shm/download_process ] ; then
	#statements
	echo "Already downloading"
	exit
fi

sudo touch /run/shm/download_process

SERVER=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/server_addr`
echo $SERVER
if [ "$SERVER"x == "x" ]; then
	#statements
	echo "No server found"
	exit
fi

LST_PACKAGE_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/lst_version`

DOWNLOAD_STATUS=true
LOCAL_IMAGES=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls config/image`
echo $LOCAL_IMAGES
for IMAGE in $LOCAL_IMAGES
	do
		CUR_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
		echo $CUR_VERSION
		LST_VERSION=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/cur_version`
		echo $LST_VERSION
		if [ "$CUR_VERSION"x != "$LST_VERSION"x ]; then
			#statements
			echo "upgrading..."
			LST_IMAGE=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/image`
			echo $LST_IMAGE

			# timeout 20m sudo docker pull $LST_IMAGE
			sudo docker pull $LST_IMAGE

			pullStatus=`sudo docker images |awk '{ print $1":"$2 }' |grep $LST_IMAGE`
			echo $pullStatus
			if [ "$pullStatus"x == "x" ]; then
				#statements
				DOWNLOAD_STATUS=false
				continue
			fi

			etcdctl --timeout '3s' --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/image "$LST_IMAGE"
			LST_SCRIPT=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/script`
			echo $LST_SCRIPT
			etcdctl --timeout '3s' --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/script "$LST_SCRIPT"

			NAME=${IMAGE#*image/}
			sudo docker tag  "$LST_IMAGE" "192.168.5.46:5000/"${NAME}:"$LST_VERSION"
		fi
	done

if [ "$DOWNLOAD_STATUS" = true ]; then
	#statements
	etcdctl --timeout '3s' --peers 127.0.0.1:4001 set config/global/download "DONE"
	etcdctl --timeout '3s' --peers 127.0.0.1:4001 set config/global/dwn_version "$LST_PACKAGE_VERSION"
fi

sudo rm /run/shm/download_process
