## flowservice功能简介   -------- 持续更新中


* T2的flow service融合了T1的eventserver、flowservice和nanomsg2nsq三者的功能；
* 取消了三者之间数据传输的topic，改为channel进行数据传输
* 通过ipc接收libra发来的事件，对无状态事件直接进行储存
* 对于有状态的事件经由streamtools进行分析判断，确定事件类型后再进行储存
* 视频的切片也由libra通过ipc发送，接收后直接进行储存
* 数据储存用的weedfs


## flowservice对外API


#### 返回数据格式
```
status：200

response body ： {'code':0,'msg':'Success','redirect':'','data':''}

对应的go语言结构
type ResponseStruct struct {

    Code int `json:"code"`

    Msg string `json:"msg"`

    Redirect string `json:"redirect"`

    Data interface{} `json:"data"`

}
```



### 1、获取事件的基本配置

* Path:/api/eventrule
* Method: GET
* data example
```
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": {
        "approaching_legacy": {
            "Enabled": true,
            "Id": "approaching_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 0.9
        },
        "distance_legacy": {
            "Enabled": true,
            "Id": "distance_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 100000
        },
        "dwellingtime_legacy": {
            "Enabled": true,
            "Id": "dwellingtime_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 1800
        },
        "population_legacy": {
            "Enabled": true,
            "Id": "population_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 40
        },
        "singlepeoplein_legacy": {
            "Enabled": false,
            "Id": "singlepeoplein_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 10,
            "WeekdayRange": 0
        },
        "velocity_legacy": {
            "Enabled": true,
            "Id": "velocity_legacy",
            "LowerBound": 0,
            "TimeRange": [
                0,
                0
            ],
            "UpperBound": 2200
        }
    }
}
```


### 2、设置事件的基本配置

* Path:/api/eventrule
* Method: PUT
* data example
```
请求数据
{
    "approaching_legacy": {
        "Enabled": true,
        "Id": "approaching_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 0.9
    },
    "distance_legacy": {
        "Enabled": true,
        "Id": "distance_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 100000
    },
    "dwellingtime_legacy": {
        "Enabled": true,
        "Id": "dwellingtime_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 1800
    },
    "population_legacy": {
        "Enabled": true,
        "Id": "population_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 40
    },
    "singlepeoplein_legacy": {
        "Enabled": false,
        "Id": "singlepeoplein_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 10,
        "WeekdayRange": 0
    },
    "velocity_legacy": {
        "Enabled": true,
        "Id": "velocity_legacy",
        "LowerBound": 0,
        "TimeRange": [
            0,
            0
        ],
        "UpperBound": 2200
    }
}
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": null
}
```


### 3、获取flowservice基本配置信息
* Path:/api/basicconfig
* Method:GET
* data example
```
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": {
        "full_video_storage_ttl_days": 3,
        "pub_db_url": "127.0.0.1:8880/api/db",
        "pub_vibo_url": "127.0.0.1:8881/api/vibo",
        "tss_ttl_days": 35
    }
}
```

```
配置项说明
"/config/libra/data/enable_color_tracking": value           //颜色粒子算法开启状态
"/config/eventserver/full_video_storage_ttl_days": value    //完整视频存储天数
"/config/eventserver/tss_ttl_days": value	                //事件存储天数
"/config/eventserver/pub_vibo_url": value                   //实时事件二次开发地址
"/config/eventserver/pub_db_url": value	                    //事件数据库二次开发地址
"/config/libra/sensor/depth_MoG_factor": value              //背景模型抗干扰程度
```


### 4、设置flowservice基本配置信息
* Path:/api/basicconfig
* Method:PUT
* data example
```
请求数据
{
    "full_video_storage_ttl_days": 3,
    "pub_db_url": "127.0.0.1:8880/api/db",
    "pub_vibo_url": "127.0.0.1:8881/api/vibo",
    "tss_ttl_days": 35
}
返回数据

```

### 5、获取区域/越线信息
* Path:/api/hotspot?name=(door/fence)
* Method:GET
* data example
```
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": {
        "door_3": {
            "Enabled": true,
            "Id": "door_3",
            "InwardAlert": true,
            "OutwardAlert": false,
            "PopulationThreshold": [],
            "RawPoints": [
                -1101.6929931640625,
                4197.98583984375,
                -29.1136474609375,
                609.5377197265625,
                1832.6214599609375,
                1043.6983642578125,
                915.5974731445312,
                4883.6455078125,
                612.6220703125
            ],
            "Type": 0,
            "cmd": false
        }
    }
}
```

### 6、设置区域/越线信息
* Parh:/api/hotspot?name=(door/fence)
* Method:PUT
* data example
```
请求数据
{
    "door_3": {
        "Enabled": true,
        "Id": "door_3",
        "InwardAlert": true,
        "OutwardAlert": false,
        "PopulationThreshold": [],
        "RawPoints": [
            -1101.6929931640625,
            4197.98583984375,
            -29.1136474609375,
            609.5377197265625,
            1832.6214599609375,
            1043.6983642578125,
            915.5974731445312,
            4883.6455078125,
            612.6220703125
        ],
        "Type": 0,
        "cmd": false
    }
}
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": null
}
```


### 7、删除区域/越线信息
* Parh:/api/hotspot?name=(door/fence)
* Method:DELETE
* data example
```
请求数据(线ID)
["door_1","door_2"]
返回数据
{
    "Code": 0,
    "Msg": "success",
    "Redirect": "",
    "Data": null
}
```

