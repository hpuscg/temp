
### 参数
* ip为皓目行为分析仪ip，默认为空
* value为是否开启截图，1为开启，0为关闭，默认为0
* buffer为tmp下最大能存储图片的空间，单位为M
* interval是每张截图之间的时间间隔，单位为ms

### 工具使用方法

* ./startCutbroad.arm -ip=192.168.12.12 -value=1 -buffer=100 -interval=20000

### 数据收集
* 收集图片的设备要求必须挂在到网管服务器上
* 每天设备端收集的图片会自动发送到网管服务器/data/picture文件夹下
* 然后自动删除历史照片
* 隔一段时间需要登录网管服务器，将/data/picture的tar.gz文件打包下载下来回传公司


// 注明：startCutbroad.sh属于批量处理的脚本，供参考