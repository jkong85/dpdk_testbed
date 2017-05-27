#/bin/bash
set -x 
echo "config the dpdk envi"
echo "setup dpdk"
./set_dpdk.sh  
sleep 3
echo "clear ovs"
./clear_ovs.sh  
sleep 3
echo "start ovs"
./start_ovs.sh
sleep 3
echo "setup ovs"
./set_ovs.sh  
sleep 5
echo "create user-1 user-2 dockers"
./start_docker.sh

