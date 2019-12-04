import os


def main():
    os.system('etcdctl set /config/global/sensor_uid ""')
    os.system('etcdctl set /config/global/sensor_sn ""')
    os.system('rm /run/shm/tegra_int')
    os.system('docker load -i eventserver177.tar.gz')
    os.system('docker load -i pioneer157.tar.gz')
    os.system('reboot')


if __name__ == '__main__':
    main()