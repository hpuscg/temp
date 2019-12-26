#!/bin/bash

# rm info log file
InfoLogFiles=`find /libra/logs -name "Libra*.log.INFO*" |xargs ls -t|awk '{line[NR]=$0} END {for(i=6 ;i<=NR;i++) print line[i]}'`
for InfoLogFile in ${InfoLogFiles}
do
rm ${InfoLogFile}
done

# rm warning log file
WARNINGLogFiles=`find /libra/logs -name "Libra*.log.WARNING*" |xargs ls -t|awk '{line[NR]=$0} END {for(i=6 ;i<=NR;i++) print line[i]}'`
for WARNINGLogFile in ${WARNINGLogFiles}
do
rm ${WARNINGLogFile}
done

# rm error log file
ERRORLogFiles=`find /libra/logs -name "Libra*.log.ERROR*" |xargs ls -t|awk '{line[NR]=$0} END {for(i=6 ;i<=NR;i++) print line[i]}'`
for ERRORLogFile in ${ERRORLogFiles}
do
rm ${ERRORLogFile}
done

# rm proc moinitor log file
ProcLogFiles=`find /libra/logs -name "Libra.proc_moinitor*" |xargs ls -t|awk '{line[NR]=$0} END {for(i=6 ;i<=NR;i++) print line[i]}'`
for ProcLogFile in ${ProcLogFiles}
do
rm ${ProcLogFile}
done