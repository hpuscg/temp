#!/bin/bash
IP="127.0.0.1"
if [ ! -f "differCompany.config" ];then
    echo "No config file!"
else
    for line in $(cat differCompany.config)
    do
        KEY=`echo ${line} | cut -d = -f 1`
        VALUE=`echo ${line} | cut -d = -f 2`

        case ${KEY} in
            "company_name")
                echo ${KEY}===${VALUE}
                curl -L http://${IP}:4001/v2/keys/config/global/${KEY} -XPUT -d value=${VALUE}
                ;;
            "level")
                echo ${KEY}===${VALUE}
                curl -L http://${IP}:4001/v2/keys/config/global/${KEY} -XPUT -d value=${VALUE}
                ;;
            "disableevent")
                echo ${KEY}===${VALUE}
                curl -L http://${IP}:4001/v2/keys/config/eventserver/${KEY} -XPUT -d value=${VALUE}
                ;;
            *)
                echo unkonw${KEY}----${VALUE}
                ;;
        esac
    done
fi
