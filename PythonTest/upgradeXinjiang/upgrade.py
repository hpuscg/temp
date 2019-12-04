
import os


IPs = [ "10.147.15.172", "10.147.17.3", "10.147.16.228", "10.147.17.91", "10.147.17.180", "10.147.17.181", "10.147.13.211",
        "10.147.13.218", "10.147.9.150", "10.147.12.214", "10.147.12.215", "10.147.16.186", "10.147.9.90", "10.147.15.118",
        "10.147.15.119", "10.147.17.217", "10.147.19.200", "10.147.11.17", "10.147.8.153", "10.147.8.155", "10.147.16.8",
        "10.147.14.233", "10.147.14.239", "10.147.8.51", "10.147.8.52"]


def main():
    for IP in IPs:
        os.system("scp -i id_rsa -r eventserver177.tar.gz pioneer157.tar.gz shellfile.py Release ubuntu@" + IP + ":/home/ubuntu")
        os.system("scp -i id_rsa reload_sensor.sh root@" + IP + ":/data/shell/_usrbin/")
        os.system("scp -i id_rsa reload_sensor.sh root@" + IP + ":/usr/bin/")
        os.system("scp -i id_rsa 70-persistent-net.rules root@" + IP + ":/etc/udev/rules.d/")
        os.system("scp -i id_rsa -r start_tf_service.sh stop_tf_service.sh ubuntu@" + IP + ":/data/shell/service/hookshell/")
        print(IP)


if __name__ == '__main__':
    main()