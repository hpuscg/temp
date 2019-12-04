import os
import time


def main():
    i = 0
    while True:
        tmp = os.popen("ps -a |grep main.go |wc -l").readline()
        int_tmp = int(tmp.strip())
        if int_tmp != 3:
            print("tmp is: ", int_tmp)
            os.system("go run main.go")
        time.sleep(1)
        print(int_tmp)
        i += 1


if __name__ == '__main__':
    main()
