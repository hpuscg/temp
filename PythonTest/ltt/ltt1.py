
'''
for ip_data in self.data:
    timestamp = ip_data['_source']['@timestamp']
    ip_temp = ip_data['_source']
    if "src_ip" not in ip_temp:
        print(ip_temp)
        print("------------")
        print(type(ip_temp))
        message_temp = ip_temp['message']
        print(type(message_temp))
        print("============")
        print("message_temp is :%s", message_temp)
        ip_temp = eval(message_temp)
    ip = ip_temp['src_ip']
    ip_timestamp_data = [ip, timestamp]
    self.list.append(ip_timestamp_data)
'''
