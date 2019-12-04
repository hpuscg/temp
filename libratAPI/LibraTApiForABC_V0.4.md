
## ======下面接口的访问端口均为8008=======
## ======目前测试机为：192.168.4.2=======

## 1、获取公司名称
method:get
url:/api/company/name
response:
===success
{
    "company_name":"DG"
## ==============================
    北京格灵深瞳信息技术有限公司:DG
    讯美科技股份有限公司:XM
    上海集光安防科技股份有限公司:JG
    深圳丽泽智能科技有限公司:LZ
    深圳市中航深亚科技有限公司:SY
    汉邦高科：HB
## ==============================
}
===fail
{
    "RespCode":500
    "RespData":err.Error()
}


## 2、获取设备的高低配置
method:get
url:/api/libra/level
response:
{
    "libra_level":"HIGH"
## ==============================
    高配版：HIGH
    低配版：LOW
## ==============================
}
===fail
{
    "RespCode":500
    "RespData":err.Error()
}

