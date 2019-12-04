import os

# 判断文件是否存在，清除已存在的config文件
if os.path.exists("/libra/judicial/differCompany.config"):
    os.remove("/libra/judicial/differCompany.config")
    os.system('touch /libra/judicial/differCompany.config"')

# 添加配置信息
companyName = input("In: 请输入公司简称：")
with open("/libra/judicial/differCompany.config", "a+") as f:
    f.write('company_name=' + companyName + '\n')
print("Out: 公司简称为：%s" % companyName)


level = input("In: 请输入设备的高低配置：")
with open("/libra/judicial/differCompany.config", "a+") as f:
    f.write('level=' + level + '\n')
print("Out: 设备配置水平为：%s" % level)


disableEvent = input("In: 请输入禁止的事件：")
with open("/libra/judicial/differCompany.config", "a+") as f:
    f.write('disableevent=' + disableEvent + '\n')
print("Out: 禁止的事件类型为：%s" % disableEvent)

# 将uid置空
os.system('etcdctl set /config/global/sensor_uid ""')

# 重启pioneer服务
os.system("reload_sensor.sh")
