
'''
port1 = 22
host2 = "192.168.7.122"
passWd = "ubuntu"

userName = "ubuntu"

ssh = paramiko.SSHClient()

ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

ssh.connect(host2, port1, userName, passWd)

stdin, stdout, stderr = ssh.exec_command("df -h")


print(stdout.read().decode())

ssh.close()
'''

'''
trans = paramiko.Transport(("192.168.7.122", 22))

trans.connect(username="ubuntu", password="ubuntu")

sftp = paramiko.SFTPClient.from_transport(trans)

sftp.put(localpath='/Users/hpu_scg/Desktop/1.txt', remotepath='/home/ubuntu/1.txt')

trans.close()
'''



