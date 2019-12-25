#!/bin/sh

pid=`ps -fu $USER |grep go |grep main |awk '{print $2}'`

if [ "$pid" != "" ]
then
    kill -15 $pid
fi

