#!/bin/bash

# IP=$(hostname -I | cut -f1 -d' ')
IP="127.0.0.1"

declare -A array ETCD=(
    # ["config/global/cache"]="internal"
    ["config/global/redis"]=""
    ["config/global/redispassword"]=""
    ["config/global/storage"]="weedfs"
    ["config/global/weedserver"]="127.0.0.1:9333"
    ["config/global/s3region"]=""
    ["config/global/s3accesskey"]=""
    ["config/global/s3secretkey"]=""
#
    ["config/topic/ca_topic"]="ca_events"
    ["config/topic/db_topic"]="db_events"
    ["config/topic/event_topic"]="events"
    ["config/topic/tss_topic"]="tss_events"
    ["config/topic/vibo_topic"]="vibo_events"
    ["config/topic/video_topic"]="slices"

    ["config/nsqaddr/ca_topic_server"]=$IP:4150
    ["config/nsqaddr/db_topic_server"]=$IP:4150
    ["config/nsqaddr/event_topic_server"]=$IP:4150
    ["config/nsqaddr/tss_topic_server"]=$IP:4150
    ["config/nsqaddr/vibo_topic_server"]=$IP:4150
    ["config/nsqaddr/video_topic_server"]=$IP:4150

    ####
    ["config/nsqaddr/prime_ca_topic_server"]=""
    ["config/nsqaddr/prime_db_topic_server"]=""
    ["config/nsqaddr/prime_event_topic_server"]=""
    ["config/nsqaddr/prime_tss_topic_server"]=""
    ["config/nsqaddr/prime_vibo_topic_server"]=""
    ["config/nsqaddr/prime_video_topic_server"]=""
    ####
#
    # ["config/eventserver/redisexpire"]="10"
    # ["config/eventserver/report"]="3600"
    ["config/eventserver/listenport"]=":1357"
    # ["config/eventserver/max_days"]="30"
    ["config/eventserver/persistence_cycle"]="60"
    ["config/eventserver/clear_cycle"]="3600"
    # ["config/eventserver/tss_internal"]="true"
    # ["config/eventserver/weedfsttl"]="1d"

    ["config/eventserver/data_path"]="/data"
    ["config/eventserver/event_status_ttl"]="10"
    ["config/eventserver/video_cache_ttl"]="100"
    ["config/eventserver/video_storage_ttl_days"]="30"
    ["config/eventserver/full_video_storage_ttl_days"]="3"
    ["config/eventserver/tss_ttl_days"]="30"

    ["config/eventserver/pub_db_url"]="127.0.0.1:8880/api/db"
    ["config/eventserver/pub_vibo_url"]="127.0.0.1:8881/api/vibo"

    # disable event type (the map in event.go file)
    ["config/eventserver/disableevent"]="[100,101,102,103,120,129,130,139,510,511,512,513,540,541,542,543,550,551,552,553,1000,1001,1002,1003]"
)

for key in "${!ETCD[@]}"
    do
        echo    $key===${ETCD[$key]}
        curl -L http://$IP:4001/v2/keys/$key -XPUT -d value=${ETCD[$key]}
    done