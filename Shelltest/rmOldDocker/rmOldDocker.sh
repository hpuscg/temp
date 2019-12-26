#!/bin/bash

# 获取设备版本
SENSOR_VERSION=`etcdctl get /config/global/cur_version`
# 判断设备版本是否以A或者R结尾，即判断设备版本信息是否是镜像版本或者全量升级版本
if [[ ${SENSOR_VERSION} == *A ]] || [[ ${SENSOR_VERSION} == *R ]]; then
    echo ${SENSOR_VERSION}
else
    echo not end of A or R
    exit
fi

# 获取设备的升级状态
DOWNLOAD_STATUS=`etcdctl get /config/global/download`
UPGRADE_STATUS=`etcdctl get /config/global/upgrade`
# 判断设备是否处于升级状态
if [[ ${DOWNLOAD_STATUS}x != "DONEx" ]] || [[ ${UPGRADE_STATUS}x != "DONEx" ]]; then
    echo sensor is upgrade
    exit
fi

# 定义repo和image_ID的字典
declare -A REPO_IDS
# 获取设备中所有的image
REPOS=`sudo docker images`
# 修改IFS
oldifs=$IFS; IFS=$'\n';
# 遍历所有的image
for REPO in ${REPOS}
do
    # 获取image的repo名字
    NAME=${REPO%% *}
    # 获取image的tag
    TMP_TAG=`eval echo ${REPO#* }`
    TAG=${TMP_TAG%% *}
    # 获取image的id
    TMP_ID=`eval echo ${TMP_TAG#* }`
    ID=`eval echo ${TMP_ID%% *}`
    # 组合image的repository
    NAME_TAG=${NAME}:${TAG}
    # 将image的repo和imageID放入字典中
    REPO_IDS[${NAME_TAG}]=${ID}
done
# 恢复IFS
IFS=${oldifs}
echo ${!REPO_IDS[*]}
echo ${REPO_IDS[*]}

# 根据etcd存储的信息删除etcd和docker的旧数据
ETCD_IMAGES=`etcdctl ls /config/image`
# 遍历etcd中存储的docker信息
for ETCD_IMAGE in ${ETCD_IMAGES}
do
    # 获取docker的名称
    ETCD_IMAGE_NAME=${ETCD_IMAGE##*/}
    # 获取该docker正在使用的版本
    IMAGE_USE_VERSION=`etcdctl get ${ETCD_IMAGE}/cur_version`
    # 获取etcd中存储的该docker的所有版本
    ETCD_IMAGE_KEYS=`etcdctl ls ${ETCD_IMAGE}`
    # 遍历etcd中存储的该docker的所有版本
    for ETCD_IMAGE_KEY in ${ETCD_IMAGE_KEYS}
    do
        # 获取etcd key的末端值，即docker的tag
        ETCD_IMAGE_VERSION=${ETCD_IMAGE_KEY##*/}
        # 判断etcd中的版本是否是当前使用版本
        if [ ${ETCD_IMAGE_VERSION} != in_use ] && [ ${ETCD_IMAGE_VERSION} != cur_version ] && [ ${ETCD_IMAGE_VERSION} != ${IMAGE_USE_VERSION} ]; then
        # 删除etcd中旧的docker的记录信息
            etcdctl rm ${ETCD_IMAGE_KEY}/image
            etcdctl rm ${ETCD_IMAGE_KEY}/script
            etcdctl rmdir ${ETCD_IMAGE_KEY}
        fi
    done
    # 当前运行的docker的name和tag
    IMAGE_USE_NAME_TAG=${ETCD_IMAGE_NAME}:${IMAGE_USE_VERSION}
    echo ${IMAGE_USE_NAME_TAG}
    # 遍历docker images
    for REPO_KEY in $(echo ${!REPO_IDS[*]})
    do
        # 判断docker image是否是当前运行的docker
        if [[ ${REPO_KEY} =~ ${ETCD_IMAGE_NAME} ]] && [[ ${REPO_KEY} != *${IMAGE_USE_NAME_TAG} ]]; then
            # 删除不是当前运行的docker
            sudo docker rmi -f ${REPO_IDS[$REPO_KEY]}
        fi
    done
done
