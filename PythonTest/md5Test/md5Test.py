#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-

import hashlib
import sys

reload(sys)
sys.setdefaultencoding('utf8')

def md5Test(pw):
    return hashlib.md5(pw).hexdigest()
    
if __name__ == "__main__":
    print(md5Test("admin123"))
