#!/bin/bash

# exe this file it reload_sensor.sh
# load service here


echo "hello" >/home/ubuntu/Songyan

    IMAGE="/config/image/armhf-libra-init"
    UPGRADE=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/upgrade_version`
    if [ "$UPGRADE"x != "x" ];then
        LST_VERSION=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/cur_version`
        LST_IMAGE=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/$LST_VERSION/image`
        LST_SCRIPT=`etcdctl --timeout '3s' --peers 127.0.0.1:4001 get $IMAGE/$LST_VERSION/script`
        UPCRIPT="./update.sh"
        LST_SCRIPT="${LST_SCRIPT} ${UPCRIPT}"
        echo "Run Libra-init at update mode: "${LST_SCRIPT/IMAGE/$LST_IMAGE} 
    	${LST_SCRIPT/IMAGE/$LST_IMAGE}		
            # run extra shell script
        if [ -d "/data/shell/extratask" ]; then
            DS=`date +%H%M%S`
            mkdir -p /tmp/$DS
            mv /data/shell/extratask/ /tmp/$DS/
            if [ -e "/tmp/${DS}/extratask/run.sh" ]; then                          
                sh /tmp/$DS/extratask/run.sh
	        echo "Run ExtraTask shell"
	    fi
	    rm -rf /tmp/$DS
        fi
    		
    
        etcdctl --timeout '3s' --peers 127.0.0.1:4001 rm $IMAGE/upgrade_version
        
        reboot
    fi


exit
