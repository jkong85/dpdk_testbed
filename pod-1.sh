#!/bin/bash
set -x
function myecho {
   var="\e[1;31m + $* + \e[0m"
   echo -e $var
}

function check-pods-created() {
  local attempt=0
  sleep 5
  while [[ ! -z "$(kubectl get pods -o wide|grep 'ContainerCreating')" ]]; do
    if (( attempt > 120 )); then
      echo "timeout waiting for creating pods" >> ~/kube/err.log
    fi
    echo "waiting for creating pods"
    attempt=$((attempt+1))
    sleep 5
  done
}

echo -e "\e[1;31m Deleting pods according yaml...\e[0m"

kubectl delete -f k8s-isolation-1.yaml
kubectl get pods -o wide
echo -e "\e[1;31m Creating pods according yaml...\e[0m"
kubectl create -f k8s-isolation-1.yaml
echo wait 5 seconds...
sleep 5
check-pods-created

echo -e "\e[1;31m Get a IPs of pods...\e[0m"
all=$(kubectl get pods -o wide |grep Running | grep dpdk-pod-1)
all="
NAME                READY     STATUS    RESTARTS   AGE       IP           NODE
isolation-1-hfkrv   1/1       Running   0          1m        10.168.0.3   192.168.56.21
"
echo "$all"
mypod=$(echo "${all}" | awk '{print $1}')
mypods=($mypod)
ips=$(echo "${all}"|awk '{print $6}')
hip=$(echo "${all}"|awk '{print $7}')
pid=$(echo "${all}"|awk 'NR==1 {print $1}')
i=0
hips=($hip)
echo "*****************************************"  $hips

for ip in $ips; do
   kubectl describe pod ${mypods[0]} |grep 'Labels:'

   echo $ip "${hips[$i]}" 
   if [ ${hips[$i]} == ${hips[0]} ] ; then
      echo -e "\e[1;31m Pod $pid ping  $ip on the SAME node ${hips[$i]} ${hips[0]} ...\e[0m"
   else
      echo -e "\e[1;31m Pod $pid ping  $ip on the DIFFERENT node ${hips[$i]} ${hips[0]}...\e[0m"
   fi
   kubectl exec $pid -it -- ping -c 4 $ip
   ((i=i+1))
done

