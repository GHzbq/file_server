#!/bin/bash

case $1 in 
	start)
		nohup ./file_server 2>&1 > /dev/null 2>&1 &
		echo "服务已启动..."
		sleep 1
	;;
	stop)
		killall file_server 
		echo "服务已停止..."
		sleep 1
	;;
	restart)
		killall file_server 
		sleep 1
		nohup ./file_server 2>&1 > /dev/null 2>&1 &
		echo "服务已重启..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac
