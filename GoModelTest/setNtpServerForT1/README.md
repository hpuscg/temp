* 主要文件包括：setNtp、ip.txt、README.md、setNtp.exe
* setNtp支持在ubuntu（64位）上运行
* setNtp.exe支持在Windows（64位）上运行
* 使用方法：
```
    * setNtp、ip.txt、setNtp.exe放在同一文件夹下
    * 修改ip.txt文件内容为要修改ntp设备的IP，格式和原有数据格式保持一致
    * 执行程序 ./setNtp -ntpServer="Ntpaddr"
    * 注意上面一步的Ntpaddr替换成实际的ntp地址
    * 将log打包返回给开发进行分析
```