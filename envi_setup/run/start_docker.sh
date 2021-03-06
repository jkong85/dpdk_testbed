#!/bin/bash

docker ps -a | grep user- | awk '{print $1}' | xargs -I {} docker rm {}

docker run -d -i -t --name user-1 -v /usr/local/var/run/openvswitch/vhost-user-1:/var/run/usvhost -v /dev/hugepages:/dev/hugepages -v /usr/src/dpdk-stable-16.11.1/x86_64-native-linuxapp-gcc:/home/dpdk/x86_64-native-linuxapp-gcc mydpdk

#docker run -d -i -t --name user-2 -v /usr/local/var/run/openvswitch/vhost-user-2:/var/run/usvhost -v /dev/hugepages:/dev/hugepages -v /usr/src/dpdk-stable-16.11.1/x86_64-native-linuxapp-gcc:/home/dpdk/x86_64-native-linuxapp-gcc mydpdk

#testpmd -l 0-1 -n 4 -m 1024 --no-pci --vdev=virtio_user0,Path=/var/run/usvhost -- -i --txqflags=0xf00 --disable-hw-vlan

#testpmd -l 0-1 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -a --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac

#testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -t --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0,2A:4E:5D:13:40:CD
