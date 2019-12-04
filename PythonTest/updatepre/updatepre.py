
import os


ips = ["", ""]
def main():
    for ip in ips:
        os.system("./upinit.sh " + ip)

if __name__ == '__main__':
    main()
