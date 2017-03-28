# dpdk_testbed

Assume vm1 and vm2, traffic starts from the docker on vm1 to the docker on vm2
1. start docker on vm1 with option <--i -a>, it will save docker's mac address to the file /home/mac.txt, and then write it to the ETCD server. docker is running on default fwd mode
2. vm2 will read the mac from etcd, and start docker on vm2 with option <--i -t  ... --eth_peer_mac ... > under the txonly fwd mode


config.go:          read the mac address from mac.txt and then write to the ETCD server
start_container.go: read the mac address from ETCD server and then start docker with the peer_mac_address
