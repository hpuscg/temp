from kafka import KafkaConsumer


def main():
    '''
    主函数
    :return:
    '''
    consumer = KafkaConsumer('test', bootstrp_servers=['192.168.7.122:9092'])

    for message in consumer:
        print("%s:%d:%d:key=%s value=%s" % (message.topic, message.partition,
                                            message.offset, message.key, message.value))
