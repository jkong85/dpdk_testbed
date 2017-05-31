#!/bin/bash

cp -rf ../test-pmd /usr/src/dpdk-stable-16.11.1/app/test-pmd

pushd $DPDK_DIR
make install T=$DPDK_TARGET DESTDIR=install
popd
