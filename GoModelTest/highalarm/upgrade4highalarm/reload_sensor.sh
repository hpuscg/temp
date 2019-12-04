#!/bin/bash

## TODO
# make sure that docker is working
sudo service docker start

C_COUNT=`sudo docker ps -a | wc -l`

if [ $C_COUNT -ne 1 ]; then
	#statements
	sudo docker rm -f $(sudo docker ps -a -q)
fi

[ -d "/data" ] || sudo mkdir /data
[ -d "/data/etcd" ] || sudo mkdir /data/etcd
[ -d "/data/nsq" ] || sudo mkdir /data/nsq
[ -d "/data/tf" ] || sudo mkdir /data/tf
[ -d "/data/vodserver" ] || sudo mkdir /data/vodserver
[ -d "/data/tf/eventserver" ] || sudo mkdir /data/tf/eventserver
# [ -d "/data/tf/weedmaster" ] || sudo mkdir /data/tf/weedmaster
# [ -d "/data/tf/weedvolume" ] || sudo mkdir /data/tf/weedvolume
[ -d "/libra" ] || sudo mkdir /libra
[ -d "/libra/judicial" ] || sudo mkdir /libra/judicial

DEFAULT_ETCD="192.168.5.46:5000/armhf-etcd:2.2.2"
PIONEER_IMAGE="192.168.5.46:5000/armhf-pioneer:1.5.4_high_alarm"
LIBRA_INIT_IMAGE="192.168.5.46:5000/armhf-libra-init:1.1.0"

ETCD=`sudo docker ps -a | awk '{ print $NF }' | grep etcd`
if [ "$ETCD"x != "etcd"x ]; then
	#statements
	echo "------ etcd ------"
	sudo docker run --net="host" --restart=always \
		-i -t -d --name etcd \
		-v /data/etcd:/etcd -v /etc/localtime:/etc/localtime:ro \
		$DEFAULT_ETCD \
		etcd -data-dir /etcd/dg -name dg \
		-advertise-client-urls http://127.0.0.1:4001 \
		-listen-client-urls http://127.0.0.1:4001 \
		-initial-advertise-peer-urls http://0.0.0.0:2380 \
		-listen-peer-urls http://0.0.0.0:2380 \
		-initial-cluster-token etcd-cluster \
		-initial-cluster dg=http://0.0.0.0:2380 \
		-initial-cluster-state new
		# --memory="200m" \
fi

sleep 1s

# 确定etcd正常启动后再执行后续步骤
# version=""
# until [[ "$version" =~ "etcd" ]]; do
# 	#statements
# 	version=`curl -L http://127.0.0.1:4001/version`
# 	echo "$version"
# done
#
etcd_status=false
for (( c=1; c<=10; c++ ))
do
	echo "check etcd for $c times"
	version=`curl -L http://127.0.0.1:4001/version`
	if [[ "$version" =~ "etcd" ]]; then
   		#statements
   		etcd_status=true
   		break
    fi
    sleep 1s
done

if [ "$etcd_status" = false ]; then
	#statements
	# TODO
	# open gate
    logger "sleep 3m before reboot because of etcd error."
	sleep 3m

	sudo reboot
fi

# 开机启动，同步时间
synctime.sh

# TODO
# seal gate


sensor_uid=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/sensor_uid`
echo $sensor_uid

if [ "$sensor_uid"x == "x" ]; then
	#statements

	echo "------ pioneer ------"
	sudo docker pull $PIONEER_IMAGE

	sudo docker run -i -t --rm \
		--privileged --net=host \
		$PIONEER_IMAGE \
		/bin/bash init.sh

	sleep 1s

	# [ -d /libra ] && sudo rm /libra/* -R

	echo "------ libra_init ------"
	sudo docker pull $LIBRA_INIT_IMAGE

	sudo docker run -i -t --rm \
		-v /libra:/libra \
		$LIBRA_INIT_IMAGE


	# 初始化结束，设置factory_default为1，需要上传或者覆盖密钥文件
	etcdctl --peers 127.0.0.1:4001 set config/global/factory_default "1"

fi

LIBRA_TASK_SIZE=`ls /libra/config/task.txt -alh | awk '{ print $5 }'`
if [[ $LIBRA_TASK_SIZE -lt 150 ]]; then
    #statements
    sudo docker run -i -t --rm \
		-v /libra:/libra \
		$LIBRA_INIT_IMAGE
fi


# sensor_uid仍然为空，初始化失败，删除/data/etcd后重启
sensor_uid=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/sensor_uid`
echo $sensor_uid
if [ "$sensor_uid"x == "x" ]; then
	#statements
	sudo rm /data/etcd/* -R
	# TODO
	# open gate
    logger "sleep 3m before reboot because of etcd error."
	sleep 3m

	sudo reboot
fi
# TODO
# seal gate

## modify DOCKER_OPTS according server address
# TODO
# DONE
SERVER=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/server_addr`

if [ "$SERVER"x != "x" ]; then
	#statements
	cat /etc/default/docker | while read line; do
		#statements
		if echo $line | grep -qe '^DOCKER_OPTS'; then
			#statements
			sudo touch /run/shm/docker_opts
			if [[ ! $line =~ "$SERVER" ]]; then
				#statements
				sudo sed -i "s/DOCKER_OPTS=\"/DOCKER_OPTS=\"--insecure-registry $SERVER:5000 /g" /etc/default/docker

				sudo service docker restart
			fi
		fi
	done

	if [ ! -f /run/shm/docker_opts ] ; then
		#statements
		echo "DOCKER_OPTS=\"--insecure-registry $SERVER:5000 -g /var/lib/docker -H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock --dns 8.8.8.8 --dns 8.8.4.4\"" >>/etc/default/docker

		sudo service docker restart
	fi
fi

sleep 1s

# 确定etcd正常启动后再执行后续步骤
etcd_status=false
for (( c=1; c<=10; c++ ))
do
	echo "check etcd for $c times"
	version=`curl -L http://127.0.0.1:4001/version`
	if [[ "$version" =~ "etcd" ]]; then
   		#statements
   		etcd_status=true
   		break
    fi
    sleep 1s
done

if [ "$etcd_status" = false ]; then
	#statements
	sudo reboot
fi

## sync with remote server
remote_etcd_addr=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/remote_etcd_addr`

PINGABLE_ETCD=false
if [ "$remote_etcd_addr"x != "x" ]; then
	#statements
	REMOTE_ETCD_ADDR=${remote_etcd_addr#*://}
	echo "etcdctl peers: "$REMOTE_ETCD_ADDR

	ETCD_ADDR=${REMOTE_ETCD_ADDR%:*}
	echo "etcd address: "$ETCD_ADDR

	PING_COUNT=`ping -c4 $ETCD_ADDR | grep 'received' | awk -F',' '{ print $2 }' | awk '{ print $1 }'`
	if [ $PING_COUNT > 0 ]; then
		#statements
		PINGABLE_ETCD=true
	fi
fi

# if get false, stop sync ; get true or get nil, continue
CHECK_VERSION=false
if [ "$PINGABLE_ETCD" = true ]; then
	CHECK_COUNT=`curl --request GET --url http://$SERVER:8008/api/config?key=/config/check_version | grep true | wc -l`
	if [ $CHECK_COUNT -ne 0 ]; then
		CHECK_VERSION=true
	fi
fi
if [ "$CHECK_VERSION" = false ]; then
	PINGABLE_ETCD=false
fi

if [ "$PINGABLE_ETCD" = true ]; then
	#statements
	LOCAL_IMAGES=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls config/image`
	# declare -A IMAGE_IN_USE
	IMAGE_IN_USE=/run/shm/image_in_use
	[ -f $IMAGE_IN_USE ] && sudo rm $IMAGE_IN_USE
	sudo touch /run/shm/image_in_use
	for IMAGE in $LOCAL_IMAGES
		do
			VALUE=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/in_use`
			if [ $VALUE != false ]; then
				#statements
				# IMAGE_IN_USE[$IMAGE]="true"
				echo $IMAGE,true >> $IMAGE_IN_USE
			fi
		done

	for IMAGE in $LOCAL_IMAGES
		do
			echo $IMAGE
			etcdctl --peers 127.0.0.1:4001 set $IMAGE/in_use "false"
		done

	TLS=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/global/remote_etcd_tls`

	if [ "$TLS"x == "1"x ]; then
		#statements

		REACHABLE_OPTIMUS=`curl -L http://$ETCD_ADDR:8008/api/tls/ca | grep CERTIFICATE`

		if [ "$REACHABLE_OPTIMUS"x != "x" ]; then
			#statements
			TLS_PATH=/usr/share/ca-certificates
			CA_FILE=$TLS_PATH/ca.crt
			CERT_FILE=$TLS_PATH/client.crt
			KEY_FILE=$TLS_PATH/client.key.insecure

			ls $TLS_PATH
			## 避免更换服务器后不能更新新的证书，每次启动重新下载证书
			##
			[ -f $CA_FILE ] && sudo rm $CA_FILE
			sudo touch $CA_FILE
			sudo curl -L http://$ETCD_ADDR:8008/api/tls/ca > $CA_FILE
			sudo chmod 0444 $CA_FILE
			##
			[ -f $CERT_FILE ] && sudo rm $CERT_FILE
			sudo touch $CERT_FILE
			sudo curl -L http://$ETCD_ADDR:8008/api/tls/cert > $CERT_FILE
			sudo chmod 0644 $CERT_FILE
			##
			[ -f $KEY_FILE ] && sudo rm $KEY_FILE
			sudo touch $KEY_FILE
			sudo curl -L http://$ETCD_ADDR:8008/api/tls/key > $KEY_FILE
			sudo chmod 0644 $KEY_FILE

			ls $TLS_PATH

			IMAGES=`etcdctl --timeout '10s' --ca-file $CA_FILE --cert-file $CERT_FILE --key-file $KEY_FILE --peers https://$REMOTE_ETCD_ADDR ls config/image`
			## TODO: bad tls
			#
			for IMAGE in $IMAGES
				do
				    # echo $IMAGE
					IMAGE_NAME=${IMAGE##*/}
					echo $IMAGE_NAME

					etcdctl --peers 127.0.0.1:4001 set $IMAGE/in_use "true"

					CUR_VERSION=`etcdctl --peers 127.0.0.1:4001 get $IMAGE/cur_version`
					echo $CUR_VERSION

					# LST_VERSION=`etcdctl --ca-file $CA_FILE --cert-file $CERT_FILE --key-file $KEY_FILE --peers https://$REMOTE_ETCD_ADDR get $IMAGE/cur_version`
					LST_VERSION=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/cur_version`
					echo $LST_VERSION
					if [[ $LST_VERSION == *"Failed"* ]]; then
						#statements
						continue
					fi

					# TO_UPGRADE=`etcdctl --ca-file $CA_FILE --cert-file $CERT_FILE --key-file $KEY_FILE --peers https://$REMOTE_ETCD_ADDR get $IMAGE/to_upgrade`
					TO_UPGRADE=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/to_upgrade`
					#

					if [ \( "$TO_UPGRADE"x == "true"x \) -a \( "$LST_VERSION"x != "x" \) ]; then
						#statements

						# LST_SCRIPT=`etcdctl --ca-file $CA_FILE --cert-file $CERT_FILE --key-file $KEY_FILE --peers https://$REMOTE_ETCD_ADDR get $IMAGE/$LST_VERSION/script`

						if [ "$CUR_VERSION"x != "$LST_VERSION"x ]; then
							echo "upgrade ..."
							#statements
							# LST_IMAGE=`etcdctl --ca-file $CA_FILE --cert-file $CERT_FILE --key-file $KEY_FILE --peers https://$REMOTE_ETCD_ADDR get $IMAGE/$LST_VERSION/image`
							LST_IMAGE=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/image`

							timeout 10m sudo docker pull $LST_IMAGE

							pullStatus=`sudo docker images |awk '{ print $1":"$2 }' |grep $LST_IMAGE`
							if [ "$pullStatus"x == "x" ]; then
								#statements
								continue
							fi
							etcdctl --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/image "$LST_IMAGE"

							LST_SCRIPT=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/script`
							etcdctl --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/script "$LST_SCRIPT"
							#
							etcdctl --peers 127.0.0.1:4001 set $IMAGE/cur_version "$LST_VERSION"
							# curl -L http://127.0.0.1:4001/v2/keys/$IMAGE/$LST_VERSION/image -XPUT -d value=$LST_IMAGE
							#
						fi
					fi
				done
		else
			# for IMAGE in $LOCAL_IMAGES
			# do
			# 	VALUE=${IMAGE_IN_USE[$IMAGE]}
			# 	etcdctl --peers 127.0.0.1:4001 set $IMAGE/in_use $VALUE
			# done
			cat $IMAGE_IN_USE | while read line; do
				#statements
				key=`echo $line | awk -F ''","'' '{print $1}'`
    			value=`echo $line | awk -F ''","'' '{print $2}'`
    			etcdctl --peers 127.0.0.1:4001 set $key/in_use $value
			done
		fi

	else
		IMAGES=`etcdctl --timeout '10s' --peers http://$REMOTE_ETCD_ADDR ls config/image`

		for IMAGE in $IMAGES
		do
		    # echo $IMAGE
			IMAGE_NAME=${IMAGE##*/}
			echo $IMAGE_NAME

			etcdctl --peers 127.0.0.1:4001 set $IMAGE/in_use -XPUT -d value="true"

			CUR_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
			# echo $CUR_VERSION
			TO_UPGRADE=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/to_upgrade`
			#
			LST_VERSION=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/cur_version`
			if [[ $LST_VERSION == *"Failed"* ]]; then
				#statements
				continue
			fi
			# echo $LST_VERSION

			if [ \( $TO_UPGRADE = true \) -a \( "$LST_VERSION"x != "x" \) -a \( "$CUR_VERSION"x != "$LST_VERSION"x \) ]; then
				#statements

				LST_IMAGE=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/image`
				timeout 5m sudo docker pull $LST_IMAGE

				pullStatus=`sudo docker images |awk '{ print $1":"$2 }' |grep $LST_IMAGE`
				if [ "$pullStatus"x == "x" ]; then
					#statements
					continue
				fi

				LST_SCRIPT=`curl --request GET --url http://$SERVER:8008/api/config?key=$IMAGE/$LST_VERSION/script`
				#
				etcdctl --peers 127.0.0.1:4001 set $IMAGE/cur_version "$LST_VERSION"
				etcdctl --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/image "$LST_IMAGE"
				# curl -L http://127.0.0.1:4001/v2/keys/$IMAGE/$LST_VERSION/image -XPUT -d value=$LST_IMAGE
				etcdctl --peers 127.0.0.1:4001 set $IMAGE/$LST_VERSION/script "$LST_SCRIPT"
				#
			fi
		done
	fi
fi

IMAGES=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls config/image`
echo $IMAGES

CUR_ETCD=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-etcd/cur_version`
etcd_image=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-etcd/$CUR_ETCD/image`
etcd_script=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-etcd/$CUR_ETCD/script`
# echo $etcd_image $etcd_script

if [ \( "$etcd_image"x != "x" \) -a \( "$DEFAULT_ETCD"x != "$etcd_image"x \) ]; then
	#statements
	sudo docker rm -f etcd

	# sudo docker pull $etcd_image

	${etcd_script/IMAGE/$etcd_image}
fi

sleep 1s

echo "------ nsqd ------"
CUR_NSQD=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-nsqd-live/cur_version`
nsqd_live_image=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-nsqd-live/$CUR_NSQD/image`
nsqd_live_script=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-nsqd-live/$CUR_NSQD/script`

# sudo docker pull $nsqd_live_image

${nsqd_live_script/IMAGE/$nsqd_live_image}

sleep 1s

echo "------ crtmpserver ------"
CUR_CRTMPSERVER=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-crtmpserver/cur_version`
crtmpserver_image=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-crtmpserver/$CUR_CRTMPSERVER/image`
crtmpserver_script=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get config/image/armhf-crtmpserver/$CUR_CRTMPSERVER/script`

# sudo docker pull $crtmpserver_image

${crtmpserver_script/IMAGE/$crtmpserver_image}

sleep 1s

###
IP=$(hostname -I | cut -f1 -d' ')
echo $IP

curl -L http://127.0.0.1:4001/v2/keys/config/global/localhost -XPUT -d value=$IP
curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_cmd_tcp -XPUT -d value=$IP:4150
curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_cmd_http -XPUT -d value=$IP:4151
curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_resp_tcp -XPUT -d value=$IP:4150
curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_report_tcp -XPUT -d value=$IP:4150

curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_tuner_cmd -XPUT -d value=$IP:4150
curl -L http://127.0.0.1:4001/v2/keys/config/nsqaddr/nsq_addr_tuner_resp -XPUT -d value=$IP:4150

curl -L http://127.0.0.1:4001/v2/keys/config/encoder/live_addr -XPUT -d value="rtsp://$IP/libra"

for IMAGE in $IMAGES
do
    echo $IMAGE
    if [ \( "$IMAGE"x == "/config/image/armhf-eventserver"x \) -o \( "$IMAGE"x == "/config/image/armhf-pioneer"x \) -o \( "$IMAGE"x == "/config/image/armhf-libra-init"x \) -o \( "$IMAGE"x == "/config/image/armhf-etcd"x \) -o \( "$IMAGE"x == "/config/image/armhf-nsqd-live"x \) -o \( "$IMAGE"x == "/config/image/armhf-crtmpserver"x \) -o \( "$IMAGE"x == "/config/image/armhf-vodserver"x \) ]; then
    	#statements
    	echo "Not now."
    	continue
    fi

    IN_USE=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/in_use`
    if [ $IN_USE = false ]; then
    	#statements
    	continue
    fi

	CUR_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
	CUR_IMAGE=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/$CUR_VERSION/image`
	CUR_SCRIPT=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/$CUR_VERSION/script`

	# sudo docker pull $CUR_IMAGE

	${CUR_SCRIPT/IMAGE/$CUR_IMAGE}

	sleep 1s
done

sudo docker start $(sudo docker ps -a -q)

#sudo docker restart bumble
dockerCount=`sudo docker ps -a -q | wc -l`
echo $dockerCount
if [ "$dockerCount"x == "1x" ]; then
	#statements
	echo "there is only $dockerCount service, reboot"
	sleep 60s
	sudo reboot
fi

echo "--- done ---"
