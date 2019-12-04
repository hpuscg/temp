
### 1、本机环境搭建

* pip install paramiko	注：此步骤在第一次使用的时候进行操作，后续不需要重复操作


### 2、定制化操作步骤


* 修改本机的IP为192.168.12.*（192.168.12.12除外），并与皓目一体机通过网线直连，保证皓目一体机与本机之间能ping通

* 解压压缩包：tar -zxvf auto_upgrade.tar.gz

* 进入auto_upgrade目录:cd auto_upgrade

* 执行python脚本，完成定制化工作：python3 auto_upgrade.py

* 执行完成之后，皓目一体机需要重启，重启完成之前请勿断电，以免操作失败(摄像头红色光束发出，即为重启完成)



