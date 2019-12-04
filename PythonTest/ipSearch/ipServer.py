import socket
import time

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

PORT = 9999

network = '<broadcast>'

def main():
    while True:
        s.sendto('Client broadcast message!'.encode('utf-8'), (network, PORT))
        print("=====server=======")
        time.sleep(1)

if __name__ == '__main__':
    main()
