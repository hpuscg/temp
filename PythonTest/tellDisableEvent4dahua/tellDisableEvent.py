
import paramiko
import os


class TellDisableEvent(object):
    """
    查询目标设备被禁止的事件
    """

    def __init__(self, ip):
        self.ip = ip
        self.userName = "ubuntu"
        self.keyFile = "id_rsa"
        self.port = 22

    def try_ping(self):
        while True:
            if os.system("ping -c 1 " + self.ip) == 0:
                print("this is ok!")
                break
            else:
                # 不能ping通，提示用户重新输入sensor ip
                print("Can't ping %s" % self.ip)
                self.ip = input("Please input the sensor ip again:")

    def catDisableEvent(self):
        try:
            # 获取ssh实例化对象
            client = paramiko.SSHClient()
            #
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.keyFile)
            # 连接目标主机
            client.connect(self.ip, self.port, username=self.userName, pkey=key_file, timeout=20)
            # 查询命令
            strCmd = "etcdctl get /config/eventserver/disableevent"
            # 执行Linux命令
            stdin, stdout, stderr = client.exec_command(strCmd)

            result = stdout.readline()

            if "" == result:
                print("No disableEvent!")
            else:
                print("DisableEvent is: %s" % result)
            # 关闭连接
            client.close()

        except Exception as e:
            print(e)

    def catCompanyName(self):
        try:
            # 获取ssh实例化对象
            client = paramiko.SSHClient()
            #
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.keyFile)
            # 连接目标主机
            client.connect(self.ip, self.port, username=self.userName, pkey=key_file, timeout=20)
            # 查询命令
            strCmd = "etcdctl get /config/global/company_name"
            # 执行Linux命令
            stdin, stdout, stderr = client.exec_command(strCmd)

            result = stdout.readline()

            if "" == result:
                print("No company name!")
            else:
                print("The company name is : %s" % result)
            # 关闭连接
            client.close()

        except Exception as e:
            print(e)

    def imagesTag(self):
        try:
            # 获取ssh实例化对象
            client = paramiko.SSHClient()
            #
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.keyFile)
            # 连接目标主机
            client.connect(self.ip, self.port, username="root", pkey=key_file, timeout=20)
            # 查询命令
            str_cmd = "docker images"
            # 执行Linux命令
            stdin, stdout, stderr = client.exec_command(str_cmd)

            result = stdout.readlines()
            # print("%s result is: %s" % (str_cmd, result))

            for ret in result:
                if "armhf-libra-cuda" in ret:
                    outPut = ret[35:].split()[0]
                    print("libra-cuda version is: %s" % outPut)
                elif "armhf-pioneer" in ret:
                    outPut = ret[35:].split()[0]
                    print("pioneer version is: %s" % outPut)
                elif "armhf-eventserver" in ret:
                    outPut = ret[35:].split()[0]
                    print("eventserver version is: %s" % outPut)
                elif "armhf-bumble-bee" in ret:
                    outPut = ret[35:].split()[0]
                    print("bumble-bee version is: %s" % outPut)
            # 关闭连接
            client.close()
        except Exception as e:
            print(e)

    def psCount(self):
        try:
            # 获取ssh实例化对象
            client = paramiko.SSHClient()
            #
            client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            # 获取密钥文件
            key_file = paramiko.RSAKey.from_private_key_file(self.keyFile)
            # 连接目标主机
            client.connect(self.ip, self.port, username="root", pkey=key_file, timeout=20)
            # 查询命令
            str_cmd = "docker ps | wc -l"
            # 执行Linux命令
            stdin, stdout, stderr = client.exec_command(str_cmd)

            result = stdout.readline()

            print("docker ps count is : %s" % result)
        except Exception as e:
            print(e)


def main():
    # 输入初始IP
    ip = input("Please input the sensor ip:")
    # 初始化查询禁止事件类对象
    tellDisableEvent = TellDisableEvent(ip)
    # 试探目标IP能否ping通
    tellDisableEvent.try_ping()
    # 获取目标IP的禁止事件
    tellDisableEvent.catDisableEvent()
    # 获取目标IP的公司简称
    tellDisableEvent.catCompanyName()
    # docker images tag
    tellDisableEvent.imagesTag()
    # docker ps count
    tellDisableEvent.psCount()


if __name__ == '__main__':
    main()

