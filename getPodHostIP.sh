#!/bin/bash
#set -x
input="
dpdk-pod-1-hfkrv   1/1       Running   0          1m        10.168.0.3   192.168.56.21
dpdk-pod-2-hfkrv   1/1       Running   0          1m        10.168.0.4   192.168.56.22
"
#all=$(kubectl get pods -o wide |grep Running | grep $1)
all=$(echo $input |grep Running | grep $1)

mypod=$(echo "${all}" | awk '{print $1}')
mypods=($mypod)
ips=$(echo "${all}"|awk '{print $6}')
hip=$(echo "${all}"|awk '{print $7}')
pid=$(echo "${all}"|awk 'NR==1 {print $1}')
hips=($hip)

if [ $2 = "1" ]; then  
    echo $mypod
fi  

if [ $2 = "2" ]; then  
    echo $ips  
fi 

if [ $2 = "3" ]; then  
    echo $hip  
fi 
