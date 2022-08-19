* 主要文件包括：setSensorIp、iplist.csv、README.md、setSensorIp.exe
* setSensorIp支持在ubuntu（64位）上运行
* setSensorIp.exe支持在Windows（64位）上运行
* 使用方法：
```
    * setSensorIp、iplist.csv、setSensorIp.exe放在同一文件夹下
    * 修改iplist.csv文件内容为要修改ntp设备的IP，格式和原有数据格式保持一致
    * 执行程序 ./setSensorIp
    * 注意上面一步的Ntpaddr替换成实际的ntp地址
    * 将log打包返回给开发进行分析
```