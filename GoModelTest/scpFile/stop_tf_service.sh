#!/bin/bash

# EVENTSERVER=`sudo docker ps -a | awk '{ print $NF }' | grep eventserver`
# if [ "$EVENTSERVER"x = "eventserver"x ]; then
#	#statements
#	sudo docker stop eventserver
# fi

VODSERVER=`sudo docker ps -a | awk '{ print $NF }' | grep vodserver`
if [ "$VODSERVER"x = "vodserver"x ]; then
	#statements
	sudo docker stop vodserver
fi
