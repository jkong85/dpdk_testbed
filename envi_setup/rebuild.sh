#!/bin/bash

pushd $DPDK_DIR
make install T=$DPDK_TARGET DESTDIR=install
popd
