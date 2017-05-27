#!/bin/bash
ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev

ovs-vsctl add-port br0 dpdkport0 -- set Interface dpdkport0 type=dpdk options:dpdk-devargs=0000:00:08.0
ovs-vsctl add-port br0 dpdkport1 -- set Interface dpdkport1 type=dpdk options:dpdk-devargs=0000:00:09.0

ovs-vsctl add-port br0 vhost-user-1 -- set Interface vhost-user-1 type=dpdkvhostuser  
ovs-vsctl add-port br0 vhost-user-2 -- set Interface vhost-user-2 type=dpdkvhostuser  
ovs-vsctl show
