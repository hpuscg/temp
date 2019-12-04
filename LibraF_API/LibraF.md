# API接口文档
* --- 持续更新中

## GB2016

### 获取GB2016设置

内容  | 说明  |  值
---|:--:|---:
请求方法|method|get
请求路径|URL   |IP/api/gb28181/status?RpcGbChanParam=RpcChannel0
参数   |RpcGbChanParam|RpcChannel0 / RpcChannel1

### 设置GB2016

内容  | 说明  |  值
---|:--:|---:
请求方法|method|post
请求路径|URL   |IP/api/gb28181?RpcGbChanParam=RpcChannel0
参数1  |RpcGbChanParam|RpcChannel0 / RpcChannel1
参数2   |json  |RpcGbParam

```
type RpcGbParam struct {
  Enable bool `thrift:"Enable,1" db:"Enable" json:"Enable"`
  LocalId string `thrift:"LocalId,2" db:"LocalId" json:"LocalId"`
  LocalPort int32 `thrift:"LocalPort,3" db:"LocalPort" json:"LocalPort"`
  SipServerId string `thrift:"SipServerId,4" db:"SipServerId" json:"SipServerId"`
  SipServerIp string `thrift:"SipServerIp,5" db:"SipServerIp" json:"SipServerIp"`
  SipServerPort int32 `thrift:"SipServerPort,6" db:"SipServerPort" json:"SipServerPort"`
  UserPwd string `thrift:"UserPwd,7" db:"UserPwd" json:"UserPwd"`
  AliveTime int32 `thrift:"AliveTime,8" db:"AliveTime" json:"AliveTime"`
  HeartPeriod int32 `thrift:"HeartPeriod,9" db:"HeartPeriod" json:"HeartPeriod"`
  TimeOutNum int32 `thrift:"TimeOutNum,10" db:"TimeOutNum" json:"TimeOutNum"`
  StreamIndex RpcStreamType `thrift:"StreamIndex,11" db:"StreamIndex" json:"StreamIndex"`
  VideoId string `thrift:"VideoId,12" db:"VideoId" json:"VideoId"`
  AudioId string `thrift:"AudioId,13" db:"AudioId" json:"AudioId"`
  AlarmId string `thrift:"AlarmId,14" db:"AlarmId" json:"AlarmId"`
}
```





