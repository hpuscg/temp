#!/bin/bash

DOCKERID=`sudo docker images |grep none | cut -c 64-75`
for ID in $DOCKERID:
    do
        sudo docker rmi $ID

    done





