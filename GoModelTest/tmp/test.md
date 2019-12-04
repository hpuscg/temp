# 第三方开发者文档

> 项目的根目录必须包含三个文件，`dc-worker.py`，`metadata.json`，`dc_settings.json` 或 `dc_settings.debug.json`。

### 1. dc-worker.py - 项目中使用 sdk

```
<?python
    from dc_client_sdk import SDK
    from redis import Redis


    sdk = SDK()  # 初始化 SDK，单例模式， 可以在多进程中多次初始化
    sdk.start()  # 开启项目的时候执行，只执行一次
    MongoClient(sdk.mongodb_uri)  # 使用 mongodb 缓存数据库
    Redis(sdk.redis_host)  # 使用 redis 缓存数据库

    def start(name, count, st, et):
        # ...
        sdk.heart()
# 上传数据
        sdk.upload([{"tieba_name":"奔驰","title":"北京奔驰","comment_id":124538274196,"comment":"到店有惊喜"},
                    {"tieba_name": "奔驰","title":"梅赛德斯-奔驰","comment_id":124538272387,"comment":"梅赛德斯-奔驰授权经销商4S店"}])

        # ...

    job_params = sdk.get_query_arguments()  # 获取 metadata.json 项目中的输入参数
    name = job_params.get('name')
    count = job_params.get('count')
    st = job_params.get('st')
    et = job_params.get('et')
    start(name, count, st, et)

    sdk.get_output_fields()  # 获取 metadata.json 项目中的输出参数类型
    sdk.metadata  # 获取 metadata.json 数据
sdk.mongodb_uri  # 获取 mongodb 缓存数据库连接地址、端口
    sdk.redis_host  # 获取 redis 缓存数据库地址
    sdk.redis_port  # 获取 redis 缓存数据库端口
    sdk.debug  # 判断当前是否是调试模式
    sdk.case_index  # 获取当前测试用例下标
    sdk.get_case(1)  # 获取下标为 1 的用例数据
    sdk.get_case  # 获取所有用例数据
    sdk.heart()  # 监控项目是否存活，可在项目代码的任何一个地方调用
    sdk.exit()  # 标记任务成功结束
    sdk.exit("error msg")  # 标记任务异常结束，并发送错误信息
?>
```
> sdk.upload 不可传空数据，上传数据的字段、字段类型一定要和 matedata.json 输出字段、字段类型一致。

> dc-worker.py 务必是项目主文件。

> 项目中若使用代理请自行处理，sdk不做任何操作。


### 2. metadata.json - 项目配置。输入、输出字段的格式规范，job 参数信息

```json
{
 "settings": {
  "display_name": "baidu tieba",
  "description": "Baidu Post",
  "params": [{
    "type": "string",
    "req": false,
    "display_name": "请输入要提取的贴吧首页URL",
    "name": "tieba_url",
    "multi": false
   },
{
    "type": "integer",
    "req": true,
    "name": "count",
    "display_name": "评论总量",
    "default": 1000,
    "max": 2000,
    "min": 0
   }
  ],
  "enabled": true,
  "date_range": true,
  "image_url": "https: //cn-cdn.stratifyd.cn/icon_tieba.png",
  "roles": "cx"
 },
 "input": {
  "name": "奔驰",
  "count": 100,
  "st": 1388591271000,
  "et": 1552275135000
 },
"output": {
  "tieba_name": {
   "type": "keyword"
  },
  "title": {
   "type": "text"
  },
  "comment_id": {
   "type": "long"
  },
  "comment": {
   "type": "text"
  }
 },
 "case": [{
   "kw": "奔驰",
         "count":100,
   "st": 1388591271000,
   "et": 1552275135000
  },
{
   "kw": "大众",
      "count":100,
   "st": 1388591271000,
   "et": 1388591271000
  }
 ]
}

```

### 3.1 dc_settings.json - SDK 配置

```json
{
    "token": "2f8c1b554f4844d5b74d24bf626abb22",
    "upload_limit": 2000,
    "dc_process_limit": 3,
    "mongodb_uri": "mongodb://127.0.0.1:27017",
    "dc_host_setting": {
        "base_url": "http://192.168.3.48:880",
        "action": {
            "upload": "/dataconnect_doc/up",
            "status": "/dataconnect_doc/status"
        }
    },
    "lo_host_setting": {
        "base_url": "http://192.168.20.3:8083",
        "action": {
            "start": "/sdk/?action=start",
            "finish":"/sdk/?action=finish",
            "heart": "/sdk/?action=heart",
            "error": "/sdk/?action=error"
        }
    }
}
```
### 3.2 dc_settings.debug.json - SDK debug 模式配置
```json
{
    "token": "2f8c1b554f4844d5b74d24bf626abb22",
    "project_id": "2a2eaf54fd7b46cca909072f9cd5e613",
    "id": "972169ad0a0d4fa3a1ee288b4ab11d91",
    "upload_limit": 2000,
    "process_limit": 5,
    "mongodb_uri": "mongodb://user:passwd@cache-mongodb.stratifyd:27017"
}
```

> 项目正式上线时不需要配置 SDK。
### 4. 安装 pipenv

`pip install pipenv`

### 5. 在项目中安装 sdk

***python2:***
`pipenv --two install "https://www.zhangzhipeng2023.cn/static/eco-sdk/DC_SDK-lasted-py2.py3-none-any.whl"`

***python3:***
`pipenv --three install "https://www.zhangzhipeng2023.cn/static/eco-sdk/DC_SDK-lasted-py2.py3-none-any.whl"`

### 6. 安装项目依赖
- 直接安装

     (1): `pipenv install requests` 或者 `pipenv --two install requests`

 - 手动安装

     (1): 编辑 Pipfile

        ```
        [packages]
           requests = "*"
           tornado = "<6"
           more ...
        ```

     (2): `pipenv install` 或者 `pipenv --two install`

> Pipfile 文件中 python_version 需要和当前项目的 python 版本一致。

### 8.启动项目

`--job_id` - 任务ID

`-m` -  项目 Metadata 文件

`-c` - SDK 配置文件

 正常启动

`pipenv run python dc-worker.py --job_id 1551424437138 -m metadata.json -c dc_settings.json`

使用调试模式启动

`pipenv run python dc-worker.py --job_id 1551424437138 -m metadata.json -c dc_settings.debug.json --debug`

使用调试模式启动测试用例，项目会自动调用  metadata.json 中 case字段用例

`pipenv run python dc-worker.py --job_id 1551424437138 -m metadata.json -c dc_settings.debug.json --debug --test`
