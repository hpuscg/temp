#!/bin/bash

date >> /home/ubuntu/scg.log
(echo "30 12 * * * program" ; crontab -l )| crontab