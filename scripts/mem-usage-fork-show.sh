#!/bin/bash
pids=$(pgrep -f "python env/daemon-loop.py")

#pids=$(ps aux | grep "\[python\]" | awk '{print $2}')

for pid in $pids
do
	echo PID: $pid
	echo Total:
	sudo cat /proc/$pid/smaps | grep -i Size | awk '{Total+=$2} END {print Total/1024/1024" GB"}'
	echo RSS:
	sudo cat /proc/$pid/smaps | grep -i rss | awk '{Total+=$2} END {print Total/1024" MB"}'
	echo PSS:
	sudo cat /proc/$pid/smaps | grep -i pss | awk '{Total+=$2} END {print Total/1024" MB"}'
	echo Private_Dirty
	sudo cat /proc/$pid/smaps | grep -i Private_Dirty | awk '{Total+=$2} END {print Total/1024" MB"}'
done
