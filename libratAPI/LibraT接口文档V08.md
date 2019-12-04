

## ======下面接口的访问端口均为8008=======


## 1、获取roi状态
```
method:get
url:/api/enable_roi
response:
==success
    {
        "0"  或者    "1"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 2、设置roi状态
```
method:post
url:/api/enable_roi/update
request：
    {
        "/config/libra/data/enable_roi":"0"
    }
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 3、获取roi监测区域
```
method：get
url：/api/roiboundary
response：
==success
    {
        "[7377,15000,0,1376,1797,0,-1596,1897,0,-9096,18000,0]"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 4、设置roi监测区域
### =====获取空间三维坐标，顺时针取点=====
```
method:post
url:/api/roiboundary/update
request:
    {
        "/config/libra/data/roi_boundary_position_array":"[7377,15000,0,1376,1797,0,-1596,1897,0,-9096,18000,0]"
    }
response:
==success
    {
        "OK"
    }
fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 5、获取系统信息
```
method:get
url:/api/systeminfo
response:
==success
    {
        "Hostname":"tegra-ubuntu",
        "CurTime":"Fri May 18 11:04:22 CST 2018",
        "Uptime":"0 Days 13 Hours 29 Minutes",
        "DockerInfo":
        [
            "GitCommit=4749651-dirty",
            "GoVersion=go1.4.2",
            "KernelVersion=3.10.40",
            "Os=linux",
            "Version=1.6.0",
            "ApiVersion=1.18",
            "Arch=arm"
        ]
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 6、获取设备版本
```
method:get
url:/api/version
response:
==success
    {
        DG-UNO V2.13.160802R
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 7、获取设备描述
```
method:get
url:/api/name
reponse:
==success
    {
        "大厅2"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 8、获取设备视频流地址
```
method:get
url:/api/liveurl
response:
==success
    {
        "rtsp://192.168.4.2/libra"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 9、获取系统状态（USB/硬盘/TF卡...）
```
method:get
url:/api/sensorstatus
response:
==success
    {
        "Usb":0,
        "Disk":0,
        "TfCard":0,
        "LocalNetwork":0,
        "RemoteNetwork":0,
        "Memory":0,
        "Service":0
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 10、设置设备描述
```
method:post
url:/api/name/update
request:
    {
        "sensor_desc":"大厅2号"
    }
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 11、按照服务名称获取其id
```
method:get
url:/api/container_id
request:
    {
        "name":"libra-cuda"
    }
response:
==sucess
    {
        "523dac45a2a977ac90380aed961d43a5ec7d2b55f5b036dcffeb4a8cdfc6fe83"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 12、重启指定的服务
```
method:post
url:/api/restart_container
request:
    {
        "id": "523dac45a2a977ac90380aed961d43a5ec7d2b55f5b036dcffeb4a8cdfc6fe83"
    }
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 13、重启设备
```
method:get
url:/api/reboot
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 14、获取设备IP、子网掩码、网关
```
method:get
url:/api/staticip
response:
==success
    {
        "address":"192.168.4.2",
        "netmask":"255.255.255.0",
        "gateway":"192.168.4.254"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 15、获取设备网管服务器地址
```
method:get
url:/api/server_address
response:
==success
    {
        "192.168.4.42"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 16、设置设备IP、子网掩码、网关
```
method:post
url:/api/staticip/update
request:
    {
        "address":"192.168.4.2",
        "netmask":"255.255.255.0",
        "gateway":"192.168.4.254"
    }
response:
==success
    {
       "OK"
    }
```


## 17、设置设备网管服务器地址
```
method:post
url:/api/server_address/update
request:
    {
        "server_address":"192.168.4.42"
    }
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
==not server
    {
        "RespCode":503
        "RespData":"this is not a server" or "this is a libra"
    }
```


## 18、获取设备时间设置（时间、校时方式、校时方式对应的IP地址或数值）
```
method:get
url:/api/synctime
response:
==success
    {
        "Mode":1,
        "Server":"192.168.4.42",
        "Ntp":"192.168.4.42",
        "Date":"Wed Apr 11 17:12:40 CST 2018"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 19、设置设备时间
```
method:post
url:/api/synctime/update
request:
    {
        "Mode":1,
        "Server":"192.168.4.42",
        "Ntp":"192.168.4.42",
        "Date":"Wed Apr 11 17:12:40 CST 2018"
    }
response:
==success
    {
        "RespCode":200,
        "RespData":"OK"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```
## ===================================
```
    {
        "/config/libra/data/enable_color_tracking": value		    //颜色粒子算法开启状态
        "/config/eventserver/full_video_storage_ttl_days": value	//完整视频存储天数
        "/config/eventserver/tss_ttl_days": value					//事件存储天数
        "/config/eventserver/pub_vibo_url": value					//实时事件二次开发地址
        "/config/eventserver/pub_db_url": value						//事件数据库二次开发地址
        "/config/libra/sensor/depth_MoG_factor": value				//背景模型抗干扰程度
    }
```
## ===================================


## 20、获取指定服务设置
```
method:get
url:/api/serviceconfig
request:
    {
        "keys": [
            "/config/libra/data/enable_color_tracking",
            "/config/eventserver/full_video_storage_ttl_days",
            "/config/eventserver/tss_ttl_days",
            "/config/eventserver/pub_vibo_url",
            "/config/eventserver/pub_db_url",
            "/config/libra/sensor/depth_MoG_factor"
        ]
    }
response:
==success
    {
        "/config/eventserver/full_video_storage_ttl_days":"3",
        "/config/eventserver/pub_db_url":"127.0.0.1:8880/api/db",
        "/config/eventserver/pub_vibo_url":"127.0.0.1:4151/pub?topic=abcd",
        "/config/eventserver/tss_ttl_days":"30",
        "/config/libra/data/enable_color_tracking":"1",
        "/config/libra/sensor/depth_MoG_factor":"2"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 21、设置指定服务数值
```
method:post
url:/api/serviceconfig/update
request:
    {
        "/config/libra/data/enable_color_tracking":"1",
        "/config/eventserver/full_video_storage_ttl_days":"3",
        "/config/eventserver/tss_ttl_days":"30",
        "/config/eventserver/pub_vibo_url":"127.0.0.1:8881/api/vibo",
        "/config/eventserver/pub_db_url":"127.0.0.1:8880/api/db",
        "/config/libra/sensor/depth_MoG_factor": "2"
     }
response:
==success
    ------未设置成功的数值------
    {
        key:value
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 22、获取服务信息列表
```
##=======request的请求参数和docker ps的字段相对应===========
method:get
url:/api/container/list
request:
    {
        "all":"true",
    ###########################    这几个参数不知何意，写上面的就行
        "size":"true",
        "limit":"10",
        "since":"true",
        "before":"true"
    ###########################
    }
response:
==sucess
    {
        "SessionId":"",
        "SensorId":"",
        "UserId":"",
        "RespType":0,
        "RespCode":200,
        "RespData":
        [
            {
                "Id":"93a2ca46f13b75dff1c1ba79a4dd01d3ebc14c648e7a92412c3821c8a3cf0183","Image":"192.168.5.46:5000/armhf-vodserver:1.0.3",
                "Command":"./vodserver",
                "Created":1523441099,
                "Status":"Up 4 minutes",
                "Names":["/vodserver"]
            },
            {
                "Id":"bdb33ce6b60f6ec864c75ac56709ed514270057cedb2cfe26c3d913d20a32d5f","Image":"192.168.5.46:5000/armhf-eventserver:1.7.2",
                "Command":"./eventserver.arm -etcdserver=http://127.0.0.1:4001 -mode=fat -report_interval=600 -log_dir_deepglint /tmp/ -mem=750",
                "Created":1523441095,
                "Status":"Up 4 minutes",
                "Names":["/eventserver"]
            },
            {
                "Id":"473286e51b664cdd7160aa2b45c9ed13941ae2b4a300a0c97fd86a096f4efaa4","Image":"192.168.5.46:5000/armhf-nsq:0.3.6.fixed",
                "Command":"nsq_to_nsq -destination-nsqd-tcp-address=192.168.4.42:4150 -destination-topic=vibo_events -nsqd-tcp-address=127.0.0.1:4150 -topic=vibo_events",
                "Created":1523441091,
                "Status":"Up 4 minutes",
                "Names":["/vibo2vibo"]
            },
            {
                "Id":"54b1bf76addaf21bb7e83cc019136587fdfa68464af5b4c7972fdfdf376a34dd","Image":"192.168.5.46:5000/armhf-tunerd:2.1.0",
                "Command":"./start.sh http://127.0.0.1:4001 20 /tunerd/config/waitress.conf","Created":1523441087,"Status":"Up 4 minutes",
                "Names":["/tunerd"]
            },
            {
                "Id":"42732c78f5aab51af896737d6e051dd27684ae7e3a42ed0b9fad07e64a40f826","Image":"192.168.5.46:5000/armhf-adu:1.0.0",
                "Command":"./simpleadu --host=0.0.0.0:8086",
                "Created":1523441086,
                "Status":"Up 4 minutes",
                "SizeRw":16745,
                "Names":["/adu"]
            },
            {
                "Id":"06ad7cef531db102768457b0217ccdc6b504f8e80efa2844cdfdbd8c94524f8f","Image":"192.168.5.46:5000/armhf-nanomsg2nsq:0.3.1",
                "Command":"./nanomsg2nsq --innernsqdaddr=localhost:4150",
                "Created":1523441083,
                "Status":"Up 4 minutes",
                "Names":["/nanomsg2nsq"]
            },
            {
                "Id":"da4229b26cea826f9689a61c4d87ca9a9b676c879cc7e66949d202eb4d9e8b3d","Image":"192.168.5.46:5000/armhf-libra-cuda:1.4.34",
                "Command":"./run_all.sh -m 1",
                "Created":1523441081,
                "Status":"Up 4 minutes",
                "SizeRw":8039617,
                "Names":["/libra-cuda"]
            },
            {
                "Id":"517682ba0eb9a437e1465ee97d55ffe8fccaa267d09581800fa6ed4ba8f83c30","Image":"192.168.5.46:5000/armhf-flowservice:2.2.5.5",
                "Command":"./flowservice.armv7 -embedded=true -binpath=bin/ -etcdaddr=http://127.0.0.1:4001","Created":1523441079,
                "Status":"Up 4 minutes",
                "Names":["/flowservice"]
            },
            {
                "Id":"a6437d188181cfc5c886cca57de7059bef127f8131fc58a254b966dfa43b332f","Image":"192.168.5.46:5000/armhf-vulcand:0.8.11.0803",
                "Command":"./vulcand --apiInterface=0.0.0.0 --etcd=http://127.0.0.1:4001",
                "Created":1523441078,"Status":"Up 4 minutes",
                "Names":["/vulcand"]
            }
        ]
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 23、获取事件报警分时段设置
```
method:get
url:/api/iterate_values
response:
==success
    {
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/approaching_legacy":
        "{
            \n  \"Enabled\": true,
            \n  \"Id\": \"approaching_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 0.9,
            \n  \"WeekdayRange\": 0\n
        }",
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/distance_legacy":
        "{
            \n  \"Enabled\": true,
            \n  \"Id\": \"distance_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 100000.0,
            \n  \"WeekdayRange\": 0\n
        }",
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/dwellingtime_legacy":
        "{
            \n  \"Enabled\": true,
            \n  \"Id\": \"dwellingtime_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 1800.0,
            \n  \"WeekdayRange\": 0\n
        }",
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/population_legacy":
        "{
            \n  \"Enabled\": true,
            \n  \"Id\": \"population_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 40.0,
            \n  \"WeekdayRange\": 0\n
        }",
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/singlepeoplein_legacy":
        "{
            \n  \"Enabled\": false,
            \n  \"Id\": \"singlepeoplein_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 10.0,
            \n  \"WeekdayRange\": 0\n
        }",
        "/config/eventbrain/alertrule/a1f20f44503532313300000500d9000d/velocity_legacy":
        "{
            \n  \"Enabled\": true,
            \n  \"Id\": \"velocity_legacy\",
            \n  \"TimeRange\":
            [\n    0,
            \n    0\n  ],
            \n  \"UpperBound\": 2200.0,
            \n  \"WeekdayRange\": 0\n
        }"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 24、获取事件基本配置
```
method:get
url:/api/basic_config
response:
==success
    {
        "/config/libra/data/abnormal_threshold":"3","/config/libra/data/cutboard_update_distance":"500","/config/libra/data/detection_minimum_height":"800","/config/libra/data/enable_abnormal_action_detection":"1","/config/libra/data/enable_darkness_detection":"1","/config/libra/data/enable_fall_detection":"1","/config/libra/data/enable_invalid_operation_detection":"0","/config/libra/data/enable_latch":"0","/config/libra/data/enable_lens_protection":"1","/config/libra/data/invalid_operation_time_threshold":"10"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 25、获取区域/越线设置
```
method:get
url:/api/hotspot
request:
    {
        "name":"door"   //  "name":"fence"
    }
response:
==success
    {
        "/door_1":
        "{
            \"Id\":\"door_1\",
            \"Enabled\":true,
            \"InwardAlert\":true,
            \"OutwardAlert\":false,
            \"InsideAlert\":false,
            \"OutsideAlert\":false,
            \"InsidePopulationAlert\":false,
            \"OutsidePopulationAlert\":false,
            \"Type\":0,
            \"RawPoints\":[-1100.42041015625,3187.2216796875,35.89404296875,287.54254150390625,4748.642578125,-11.589599609375,404.3072814941406,3531.75048828125,49.24609375],
            \"PopulationThreshold\":[]
        }"
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 26、设置越线或区域信息
```
method:post
url:/api/hotspot/update
request:
    {
        "data": "{\n  \"door_1\": {\n    \"Enabled\": true,\n    \"Id\": \"door_1\",\n    \"InwardAlert\": true,\n    \"OutwardAlert\": false,\n    \"PopulationThreshold\": [],\n    \"RawPoints\": [\n      -1300.42041015625,\n      3187.2216796875,\n      35.89404296875,\n      287.54254150390625,\n      4748.642578125,\n      -11.589599609375,\n      404.30728149414062,\n      3531.75048828125,\n      49.24609375\n    ],\n    \"InsideAlert\": false,\n    \"OutsideAlert\": false,\n    \"InsidePopulationAlert\": false,\n    \"OutsidePopulationAlert\": false,\n    \"Type\": 0,\n    \"cmd\": false\n  }\n}",
        "name": "door"
    }
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 27、设置基本事件数据
```
method:post
url:/api/basic_config/update
request:
    {
        "/config/libra/data/abnormal_threshold": "3",
        "/config/libra/data/cutboard_update_distance": "500",
        "/config/libra/data/detection_minimum_height": "800",
        "/config/libra/data/enable_abnormal_action_detection": "1",
        "/config/libra/data/enable_darkness_detection": "1",
        "/config/libra/data/enable_fall_detection": "1",
        "/config/libra/data/enable_invalid_operation_detection": "1",
        "/config/libra/data/enable_latch": "0",
        "/config/libra/data/enable_lens_protection": "1",
        "/config/libra/data/invalid_operation_time_threshold": "10"
    }
response:
==success
------未设置成功的数值------
    {
        key:value
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 28、删除指定的越线或区域设置
```
method:post
url:/api/hotspot/del
request:
    {
        "name":"door",
        "ids":["door_1"]
    }
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 29、获取SensorId
```
method:get
url:/api/sensorid
response:
==success:
    {
        sensorid
    }
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 30、删除指定的事件
```
method:post
url:/api/service/del
request：
    {
        "keys":["/config/eventbrain/alertrule/a1f20f44503532313300000400b90109/door_1"]
    }
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 31、获取TF卡信息
```
method:get
url:/api/tf_status
response:
==success
    {
        "Size":"59G",
        "Avial":"2.1G",
        "Status":0
    }
==fail
    {
        "RespCode":500 || 503
        "RespData":err.Error()
    }
```


## 32、格式化TF卡
```
method:post
url:/api/tf_status/update
response:
==success
    {
        "OK"
    }
==fail
    {
        "RespCode":503
        "RespData":err.Error()
    }
```


### ================================
## 网管服务器API

### ======下面接口的访问端口均为8008=======

## 1、获取网管服务器软件版本
```
method:get
url:/api/server_version
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 2、获取挂载在网管服务器下的设备列表
```
method:get
url:/api/sensor_list
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```


## 3、获取升级时网管服务器软件版本
```
method:get
url:/api/sensor_version
response:
==success
    {}
==fail
    {
        "RespCode":500
        "RespData":err.Error()
    }
```



























































