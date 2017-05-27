#!/bin/bash
ps -axu | grep ovs | awk '{print $2}' | xargs -I {} kill -9 {}
