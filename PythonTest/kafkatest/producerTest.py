from kafka import KafkaProducer


def main():
    '''
    主函数
    :return:
    '''
    producer = KafkaProducer(bootstrap_servers=['localhost:9092'])

    for i in range(3):
        msg = "msg%d" % i
        producer.send('test', msg)
    producer.close()


if __name__ == '__main__':
    main()
