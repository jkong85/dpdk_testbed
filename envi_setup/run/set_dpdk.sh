#!/bin/bash
sysctl -w vm.nr_hugepages=8
mkdir -p /dev/hugepages
mount -t hugetlbfs none /dev/hugepages

mv /usr/src/dpdk-stable-16.11.1/app/test-pmd/ /usr/src/dpdk-stable-16.11.1/app/test-pmd-bak
cp -rf ../../test-pmd /usr/src/dpdk-stable-16.11.1/app/test-pmd

pushd $DPDK_DIR
make install T=$DPDK_TARGET DESTDIR=install
modprobe uio
insmod $DPDK_BUILD/kmod/igb_uio.ko
popd

$DPDK_DIR/tools/dpdk-devbind.py --bind=igb_uio 0000:00:08.0
$DPDK_DIR/tools/dpdk-devbind.py --bind=igb_uio 0000:00:09.0

$DPDK_DIR/tools/dpdk-devbind.py --status

