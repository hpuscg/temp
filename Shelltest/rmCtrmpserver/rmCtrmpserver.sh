#! /bin/bash

NotNone=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls /config/image/armhf-crtmpserver`
if [ "${NotNone}"x == x ]; then
echo "No crtmpserver"
exit
fi
echo ${NotNone}

echo "etcdctl rm /config/image/armhf-crtmpserver/in_use"
echo "etcdctl rm /config/image/armhf-crtmpserver/cur_version"
Tags=`etcdctl --timeout '10s' --peers 127.0.0.1:4001 ls /config/image/armhf-crtmpserver`

for Tag in ${Tags}
do
if [ \( ${Tag}x == "/config/image/armhf-crtmpserver/in_use"x \) -o \( ${Tag}x == "/config/image/armhf-crtmpserver/cur_version"x \) ]; then
echo ${Tag}
continue
fi
echo ${Tag}
echo "etcdctl rm ${Tag}/image"
echo "etcdctl rm ${Tag}/script"
echo "etcdctl rmdir ${Tag}"
done
echo "etcdctl rmdir /config/image/armhf-crtmpserver"

