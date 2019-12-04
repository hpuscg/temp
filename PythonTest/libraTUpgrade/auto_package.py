import os


Docker_image = [("eventserver-1.7.7.tar.gz", "eventserver", "1.7.7", "sudo docker run -tid --name eventserver --memory=800m --restart=always --net=host -v /tmp:/tmp -v /data/tf/eventserver:/data/slice -v /data/eventserver:/data/event -v /etc/localtime:/etc/localtime:ro IMAGE ./eventserver.arm -etcdserver=http://127.0.0.1:4001 -mode=fat -report_interval=600 -log_dir_deepglint /tmp/ -mem=750"),
                ("pioneer-1.5.7.tar.gz", "pioneer", "1.5.7", "sudo docker run -i -t --rm --privileged --net=host IMAGE /bin/bash init.sh")]
Version = "V2.13.190101R"


def main():
    cmd = "tar -zcvf package.tar.gz upgrade.sh"
    with open("upgrade.sh", "w+") as f :
        f.write("#!/bin/bash\n\n")
        f.write("IP=$(hostname -I | cut -f1 -d' ')\n\n")
        for package, server, tag, script in Docker_image:
            f.write("## modify " + server + " service\n")
            f.write("sudo docker tag 192.168.5.46:5000/armhf-" + server + ":" + tag + " $IP:5000/armhf-" + server + ":" + tag + "\n")
            f.write("curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-" + server + "/cur_version -XPUT -d value=" + tag + "\n")
            f.write('curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-' + server + '/to_upgrade -XPUT -d value="true"\n')
            f.write('curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-' + server + '/' + tag + '/image -XPUT -d value="192.168.5.46:5000/armhf-' + server + ':' + tag + '"\n')
            f.write('curl -L http://127.0.0.1:2379/v2/keys/config/image/armhf-' + server + '/' + tag + '/script -XPUT -d value="' + script + '"\n\n')
            os.system("sudo docker pull 192.168.5.46:5000/armhf-" + server + ":" + tag)
            os.system('sudo docker save --output="' + package + '" 192.168.5.46:5000/armhf-' + server + ':' + tag)
            cmd = cmd + " " + package
        f.write("## modify system version 版本号\n")
        f.write('curl -L http://127.0.0.1:2379/v2/keys/config/global/cur_version -XPUT -d value="' + Version + '"\n')
        f.close()
    print(cmd)
    os.system(cmd)


if __name__ == '__main__':
    main()

