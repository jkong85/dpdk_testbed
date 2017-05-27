#!/bin/bash
echo "export DPDK_DIR=/usr/src/dpdk-stable-16.11.1" >> /etc/environment
echo "export DPDK_TARGET=x86_64-native-linuxapp-gcc" >> /etc/environment
echo "export DPDK_BUILD=/usr/src/dpdk-stable-16.11.1/x86_64-native-linuxapp-gcc" >> /etc/environment
echo "export RTE_SDK=/usr/src/dpdk-stable-16.11.1" >> /etc/environment
