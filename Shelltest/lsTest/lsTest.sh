
timestamp2=$((`date '+%s'`*1000+`date '+%N'`/1000000))
echo $timestamp2

ls /usr/bin/check_cpu_frequency.sh /usr/bin/check_docker_config.sh /usr/bin/clear_docker_images.sh /usr/bin/clear_docker_logs.sh /usr/bin/clear_system_logs.sh /usr/bin/deletelog.sh /usr/bin/load_libraT_utils.sh /usr/bin/open_the_gate.sh /usr/bin/reboot.sh /usr/bin/reload_sensor.sh /usr/bin/reload_sensor_template.sh /usr/bin/seal_the_gate.sh /usr/bin/synctime.sh /usr/bin/tegra_init.sh /usr/bin/upload_ssh_key.sh /usr/bin/version_etcd_init.sh /usr/bin/vulcand_etcd_init.sh /usr/bin/weekly_clear.sh -lh >> /home/ubuntu/${timestamp2}
