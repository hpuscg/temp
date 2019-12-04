# import json
import os
import re


# data_file_topic_list = []
page_topic_list = []


'''
def data_file_topic():
    """
    data文件中存储的topic信息
    :return:
    """
    filename = "/data/nsq/nsqd.368.dat"
    with open(filename) as f:
        line = f.readline()
        data = json.loads(line)
        topics = data["topics"]
        i = 0
        for topic in topics:
            name = topic["name"]
            print(name)
            data_file_topic_list.append(name)
            i += 1
        print("data_file_topic数量", i)
        print(data_file_topic_list)
'''


def page_topic(page_sensor_address):
    """
    网页stats访问得到的topic信息
    :return:
    """
    resp = os.popen('ping -c 4 ' + page_sensor_address + ' -t 10').read()
    ret = re.findall(r'from 192.168.4.2: icmp_seq=', resp)
    if len(ret) != 4:
        print("一体机IP地址不正确，请重新尝试！")
        return
    body = os.popen('curl http://' + page_sensor_address + ':4151/stats').read()
    result = re.findall(r'\[.+]', body)
    j = 0
    for ret in result:
        name = ret[1: len(ret) - 1].strip()
        page_topic_list.append(name)
        j += 1
    print("page_topic数量", j)
    # print(page_topic_list)
    flag = "topic"
    k = 0
    topic_name = None
    for ct_name in page_topic_list:
        if ct_name.startswith("channel"):
            if "topic" == flag:
                topic_name = page_topic_list[k - 1]
            flag = "channel"
        else:
            if "channel" == flag:
                delete_topic(topic_name, page_sensor_address)
            flag = "topic"
            topic_name = ct_name
        k += 1


'''
def differ_topic():
    """
    data文件中与网页stats中不同的topic
    :return:
    """
    count = 0
    # data文件中有而网页stats中没有的topic
    for topic in data_file_topic_list:
        if topic not in page_topic_list:
            print(topic)
            count += 1
    print(count)
    # 网页stats中有而data文件中没有的topic
    for topic in page_topic_list:
        if topic not in data_file_topic_list:
            count += 1
    print(count)
'''


def delete_topic(topic_name, topic_sensor_address):
    os.popen('curl -X POST http://' + topic_sensor_address + ':4151/topic/delete?topic=' + topic_name)


def delete_channel(topic_name, channel_name, channel_sensor_address):
    os.popen('curl -X POST http://' + channel_sensor_address + ':4151/channel/delete?topic=' + topic_name
             + '&channel=' + channel_name)


if __name__ == "__main__":
    sensor_address = input("please input your sensor ip:")
    # data_file_topic()
    
    page_topic(sensor_address)
    # differ_topic()
