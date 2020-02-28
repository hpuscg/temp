#!/bin/bash

if [ -f /run/shm/upgrade_process ] ; then
	#statements
	echo "Already upgrading"
	exit
fi

sudo touch /run/shm/upgrade_process

SERVER=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/server_addr`
echo $SERVER
if [ "$SERVER"x == "x" ]; then
	#statements
	echo "No server found"
	exit
fi

LOCAL_IMAGES=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls config/image`
echo $LOCAL_IMAGES
for IMAGE in $LOCAL_IMAGES
	do
		CUR_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
		echo $CUR_VERSION
		LST_VERSION=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/cur_version`

		if [[ $LST_VERSION == *"Failed"* ]]; then
			#statements
			continue
		fi
		echo $LST_VERSION
		if [ \( "$CUR_VERSION"x != "$LST_VERSION"x \) -a \( "$LST_VERSION"x != "x" \) ]; then
			#statements
                    etcdctl --timeout '3s' --peers 127.0.0.1:4001 set $IMAGE/cur_version "$LST_VERSION"
		    if [ "$IMAGE"x == "/config/image/armhf-libra-init"x ]; then
		        etcdctl --timeout '3s' --peers 127.0.0.1:4001 set $IMAGE/upgrade_version "$LST_VERSION"
		    fi
		fi
# modify by songyan 
	done

UPGRADE_STATUS=true
for IMAGE in $LOCAL_IMAGES
	do
		CUR_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
		echo $CUR_VERSION
		LST_VERSION=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/cur_version`
		echo $LST_VERSION
		if [[ $LST_VERSION == *"Failed"* ]]; then
			#statements
			continue
		fi
		if [ \( "$CUR_VERSION"x != "$LST_VERSION"x \) -a \( "$LST_VERSION"x != "x" \) ]; then
			#statements
			UPGRADE_STATUS=false
			break
		fi
	done

if [ "$UPGRADE_STATUS" = true ]; then
	#statements
	etcdctl --timeout '3s' --peers 127.0.0.1:4001 set config/global/upgrade "DONE"
	LST_DWN_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/dwn_version`
	etcdctl --timeout '3s' --peers 127.0.0.1:4001 set config/global/cur_version "$LST_DWN_VERSION"
	VERPREFIX="Release Version: DG-UNO"
	VERTEXT="1c ${VERPREFIX} ${LST_DWN_VERSION}"
	echo "Update Release file"
	sed -i "$VERTEXT" /home/ubuntu/Release
  
	sudo rm /run/shm/upgrade_process

	sudo reboot
	
	exit
fi

sudo rm /run/shm/upgrade_process

