#!/bin/bash
#ps -axu | grep ovs | awk '{print $2}' | xargs -I {} kill -9 {}

rm /usr/local/etc/openvswitch/conf.db                                                                                                                                                   
rm -rf /usr/local/var/run/openvswitch/*                                                                                                                                                     

cd /usr/src/openvswitch-2.7.0/                                                                                                                                                          
ovsdb-tool create /usr/local/etc/openvswitch/conf.db vswitchd/vswitch.ovsschema                                                                                                         
ovsdb-server --remote=punix:/usr/local/var/run/openvswitch/db.sock --remote=db:Open_vSwitch,Open_vSwitch,manager_options --private-key=db:Open_vSwitch,SSL,private_key --certificate=db:Open_vSwitch,SSL,certificate --bootstrap-ca-cert=db:Open_vSwitch,SSL,ca_cert --pidfile --detach --log-file                                                                                  
#start vswitchd                                                                                                                                                                         
export DB_SOCK=/usr/local/var/run/openvswitch/db.sock                                                                                                                                   
ovs-vsctl --no-wait set Open_vSwitch . other_config:dpdk-init=true                                                                                                                      
ovs-vswitchd unix:$DB_SOCK --pidfile --detach &                                                                                                                                         
ovs-vsctl --no-wait set Open_vSwitch . other_config:dpdk-socket-mem="1024,0"                                                                                                            
ovs-vsctl set Open_vSwitch . other_config:pmd-cpu-mask=0x3   
