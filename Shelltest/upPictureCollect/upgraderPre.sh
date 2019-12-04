#!/bin/bash

etcdctl set /config/global/server_addr "server_ip"
etcdctl set /config/global/remote_etcd_addr "server_ip:4001"
etcdctl set /config/global/date_mode 1
etcdctl set /config/global/ntp_addr "ntp_ip"
