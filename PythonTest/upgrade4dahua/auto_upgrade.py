import paramiko
import os
import sys
import datetime
import time


class Auto_upgrade(object):
    """
    远程自动化完成LibraT定制化作业
    """

    # 初始化参数
    def __init__(self):
        """
        初始化参数
        """
        self.ip = "192.168.12.12"  # 目标机IP
        self.username = "root"  # 目标机用户名
        self.key_file = "id_rsa"  # 密钥文件
        self.port = 22  # 端口

    # 测试能否ping通
    def try_ping(self):
        """
        测试目标是否能ping通
        :return:
        """
        # 测试能否pin通
        if os.system("ping -c 1 " + self.ip) == 0:
            print("this is ok!")
        else:
            # 不能ping通，退出程序
            print("can't ping %s, please checkout your local ip" % self.ip)
            sys.exit(1)

    # 创建远程连接
    def try_connect(self):
        # 获取ssh实例化对象
        client = paramiko.SSHClient()
        #
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        # 获取密钥文件
        key_file = paramiko.RSAKey.from_private_key_file(self.key_file)
        # 连接目标主机
        client.connect(self.ip, self.port, username=self.username, pkey=key_file, timeout=20)
        return client

    # 远程执行Linux命令
    def try_exec(self, str_cmd, connect):
        """
        通过密钥远程执行Linux命令
        :param str_cmd:
        :param connect:
        :return:
        """
        try:
            # 执行Linux命令
            stdin, stdout, stderr = connect.exec_command(str_cmd)

            result = stdout.readline()

            print("%s result is: %s" % (str_cmd, result))
        except Exception as e:
            print(e)

    # 远程传输文件
    def try_ftp(self, from_file, to_file):
        """
        通过密钥远程传输文件
        :return:
        """
        print("start copy ", from_file)
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
        print(from_file, " copy over")

    # 删除旧的docker image
    def deleteImage(self, connect, images):
        """
        删除旧的image
        :param connect:
        :return:
        """
        # 删除旧的docker image并清理etcd上的数据
        print("start delete image")
        for image, newTag in images:
            imageName = image.split(".")[0]
            stdin, stdout, stderr = connect.exec_command("docker images |grep " + imageName)
            ret = stdout.readlines()
            for res in ret:
                tag = res.split(" ", 1)[1].lstrip().split(" ", 1)[0]
                containerId = res.split(" ", 1)[1].lstrip().split(" ", 1)[1].lstrip().split(" ", 1)[0]
                print("tag is :", tag)
                print("container id is: ", containerId)
                if "" == tag or newTag == tag:
                    print("new tag: ", newTag)
                    print("tag :", tag)
                    continue
                connect.exec_command("etcdctl rm /config/image/armhf-" + image + "/" + tag + "image")
                connect.exec_command("etcdctl rm /config/image/armhf-" + image + "/" + tag + "script")
                connect.exec_command("etcdctl rmdir /config/image/armhf-" + image + "/" + tag + "/")
                connect.exec_command("docker rmi " + containerId)
                time.sleep(10)

    # 加载更新的docker image
    def loadFiles(self, connect, delimages):
        print("start load image file")
        # 固定传输的文件
        ftp_files = [("Release", "/home/ubuntu/Release"),
                     ("70-persistent-net.rules", "/etc/udev/rules.d/70-persistent-net.rules"),
                     ("reload_sensor.sh", "/data/shell/_usrbin/reload_sensor.sh"),
                     ("differCompany.config", "/libra/judicial/differCompany.config")]
        # 需要更新的docker image包
        for delImage, _ in delimages:
            tempDelImage = (delImage, "/home/ubuntu/" + delImage)
            ftp_files.append(tempDelImage)
        # 远程传输文件，并加载新的docker image
        for from_file, to_file in ftp_files:
            print(from_file, "---", to_file)
            self.try_ftp(from_file, to_file)
            if to_file.endswith(".tar"):
                connect.exec_command("docker load --input " + to_file)
                time.sleep(10)
                connect.exec_command("rm " + to_file)
        print("load files over")

    def readImage(self):
        # 设备上所有的docker image
        allImages = ["libra-init", "tunerd", "bumble-bee", "pioneer", "eventserver", "libra-cuda", "flowservice", "adu",
                     "nanomsg2nsq", "vulcand", "crtmpserver", "vodserver", "etcd", "nsq", "onvifserver"]
        # 需要删除的docker image
        deleteImages = []
        files = os.listdir("./")
        # 根据本地的tar包确定需要删除的docker image
        for file in files:
            if file.endswith(".tar"):
                fileName = file.split(".")[0]
                fileTag = file.split(fileName + ".")[1].split(".tar")[0]
                fileTar = fileName + "." + fileTag + ".tar"
                print("file tar is: ", fileTar)
                for tempFile in allImages:
                    if tempFile == fileName:
                        deleteImages.append([fileTar, fileTag])
        return deleteImages


def main():
    # 实例化LibraT远程自动化定制类
    auto_upgrade = Auto_upgrade()
    # 检测目标主机是否能ping通
    auto_upgrade.try_ping()
    # 创建远程连接
    connect = auto_upgrade.try_connect()
    # 置空uid
    auto_upgrade.try_exec('etcdctl set /config/global/sensor_uid ""', connect)
    # 置空sn号
    auto_upgrade.try_exec('etcdctl set /config/global/sensor_sn ""', connect)
    # 删除已启动标识
    auto_upgrade.try_exec('rm /run/shm/tegra_init', connect)
    # 读取需要更新的image
    deleteImage = auto_upgrade.readImage()
    # 导入docker image
    auto_upgrade.loadFiles(connect, deleteImage)
    # 删除所有在运行的container
    auto_upgrade.try_exec('docker rm -f $(docker ps -a -q)', connect)
    # 删除需要更新的image
    auto_upgrade.deleteImage(connect, deleteImage)
    # 重启设备
    auto_upgrade.try_exec("reboot", connect)


if __name__ == '__main__':
    # 起始时间
    start_time = datetime.datetime.now()
    main()
    # 终止时间
    end_time = datetime.datetime.now()
    # 程序总用时
    time_length = (end_time - start_time).seconds
    print("time length is:", time_length)
