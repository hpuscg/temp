import paramiko
import os
import sys
import datetime


class Auto_upgrade(object):
    """
    远程自动化完成LibraT定制化作业
    """

    def __init__(self):
        """
        初始化参数
        """
        self.ip = "192.168.12.12"  # 目标机IP
        self.username = "root"  # 目标机用户名
        self.key_file = "id_rsa"  # 密钥文件
        self.port = 22  # 端口

    def try_ping(self):
        """
        测试目标是否能ping通
        :return:
        """
        # 测试能否pin通
        if 0 == os.system("ping -n 1 " + self.ip):
            print("this is ok!")
        else:
            # 不能ping通，退出程序
            print("can't ping %s, please checkout your local ip" % self.ip)
            sys.exit(1)

    def try_exec(self, str_cmd):
        """
        通过密钥远程执行Linux命令
        :param str_cmd:
        :return:
        """
        try:
            # 获取ssh实例化对象
            client = paramiko.SSHClient()
            #
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.key_file)
            # 连接目标主机
            client.connect(self.ip, self.port, username=self.username, pkey=key_file, timeout=20)
            # 执行Linux命令
            stdin, stdout, stderr = client.exec_command(str_cmd)

            result = stdout.readline()

            print("%s result is: %s" % (str_cmd, result))
            # 关闭连接
            client.close()

        except Exception as e:
            print(e)

    def try_ftp(self, from_file, to_file):
        """
        通过密钥远程传输文件
        :return:
        """
        try:
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.key_file)
            # 实例化ftp对象
            trans = paramiko.Transport((self.ip, self.port))
            # 连接目标主机
            trans.connect(username=self.username, pkey=key_file)
            # 创建传输通道
            sftp = paramiko.SFTPClient.from_transport(trans)
            # 远程传输文件
            sftp.put(localpath=from_file, remotepath=to_file)
            # 关闭连接
            trans.close()
        except Exception as e:
            print(e)


def main():
    # 实例化远程自动化LibraT定制化类
    auto_upgrade = Auto_upgrade()
    # 检测目标主机是否能ping通
    auto_upgrade.try_ping()
    # 置空sn号和uid
    auto_upgrade.try_exec('etcdctl set /config/global/sensor_uid ""')
    auto_upgrade.try_exec('etcdctl set /config/global/sensor_sn ""')

    # 删除已启动标记
    auto_upgrade.try_exec('rm /run/shm/tegra_init')
    # 初始化文件参数
    ftp_files = [("bumble.tar", "bumble.tar"), ("eventserver.tar", "eventserver.tar"),
                 ("libra.tar", "libra.tar"), ("pioneer.tar", "pioneer.tar"), ("Release", "/home/ubuntu/Release"),
                 ("70-persistent-net.rules", "/etc/udev/rules.d/70-persistent-net.rules"),
                 ("reload_sensor.sh", "/data/shell/_usrbin/reload_sensor.sh")]
    # 拷贝文件
    auto_upgrade.try_exec('mkdir /libra/judicial')
    for from_file, to_file in ftp_files:
        print("begin copy " + from_file)
        auto_upgrade.try_ftp(from_file, to_file)
        if to_file.endswith(".tar"):
            auto_upgrade.try_exec("docker load --input " + to_file)
            auto_upgrade.try_exec("rm " + to_file)
    # 重启设备
    print("copy file over")
    # 初始化需要删除的image
    images = [("pioneer", "1.5.4"), ("eventserver", "1.7.2"), ("bumble-bee", "1.8.0"), ("libra-cuda", "1.4.34")]
    # 删除所有的container
    auto_upgrade.try_exec('docker rm -f $(docker ps -a -q)')
    print("begin delete images")
    for image, version in images:
        # 删除image
        auto_upgrade.try_exec('etcdctl rm /config/image/armhf-' + image + '/' + version + '/image')
        # 删除script
        auto_upgrade.try_exec('etcdctl rm /config/image/armhf-' + image + '/' + version + '/script')
        # 删除节点
        auto_upgrade.try_exec('etcdctl rmdir /config/image/armhf-' + image + '/' + version)
        # 删除docker image
        if "eventserver" == image or "libra-cuda" == image or "pioneer" == image:
            continue
        auto_upgrade.try_exec('docker rmi $(docker images |grep ' + version + '|cut -c59-70)')
    auto_upgrade.try_exec('reboot')


if __name__ == '__main__':
    start_time = datetime.datetime.now()
    main()
    end_time = datetime.datetime.now()
    time_length = (end_time - start_time).seconds
    print(time_length)
    print("please wait about five minutes")
