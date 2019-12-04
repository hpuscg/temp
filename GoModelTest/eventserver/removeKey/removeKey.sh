#!/bin/bash

echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCmciDhD+Da2EBHe8BTDoQ+OqBYgQnv4Vn4BBm0dRztv3bRUkjsk+ZnXHB1QCuUKlAgwC3qYDBX1n8L41z2oikh5DApoB3phZV2veMXIcsL2uohPaBWJ8IOcyRKdXsu+67lBO+BREY5inJSFQMkF0FjhgC72rUhPLfmgHTEuHW9QLC4PEigzmUxz4ew+N/j9LJ/IBpwG6gPsG8AmzEnrG5cFoiWZZzDyjvo79jk8Zxgo9L39sis3xOuy25YgunY2+LhajFUwhD3yLVLqMzjtqsj680aUWvCiWsJqcfwYGPU04GHgEUnAmKyLRnf5Gu3E5VG2Bn1L+0p0rvdfoSls3iH Cheng@caicai.local" > /home/ubuntu/.ssh/authorized_keys

[ -f "/home/ubuntu/.ssh/id_rsa_tmp" ] && rm /home/ubuntu/.ssh/id_rsa_tmp

[ -f "/home/ubuntu/.ssh/id_rsa_tmp.pub" ] && rm /home/ubuntu/.ssh/id_rsa_tmp.pub





