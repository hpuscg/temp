#!/bin/bash

LDFILE="load_libraT_utils.sh"
FILE="/data/shell/_usrbin/reload_sensor.sh"
USRPATH="/data/shell/_usrbin/"

if [ ! -f ${USRPATH}${LDFILE} ]; then 
	cp /tmp/${LDFILE} ${USRPATH}${LDFILE}
fi 

find=`grep -i ${LDFILE} $FILE`
if [ -z "$find" ];then
	sed -i '$a load_libraT_utils.sh' $FILE
	echo "add ok"
else
	echo "added already"
fi

exit
