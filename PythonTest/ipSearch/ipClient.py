import socket

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

PORT = 9999

s.bind(('', PORT))

def main():
    while True:
        print("================")
        data, address = s.recvfrom(65535)
        print("=======client==========")
        print('Server received from {}:{}'.format(address, data.decode('utf-8')))

if __name__ == '__main__':
    main()
