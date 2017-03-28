package main

import (
	"fmt"
	//"io/ioutil"
	"github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"golang.org/x/net/context"
	"log"
	"os/exec"
	"time"
)

func main() {

	etcd_server := "127.0.0.1:4001"

	etcd.Register()

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{etcd_server},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		log.Fatal("Cannot create store")
	}

	key := "/10.145.240.185"

	stopCh := make(<-chan struct{})

	events, err := kv.Watch(key, stopCh)

	for {
		select {
		case pair := <-events:

			proto := "http://" + etcd_server
			cfg := client.Config{
				Endpoints:               []string{proto},
				Transport:               client.DefaultTransport,
				HeaderTimeoutPerRequest: time.Second,
			}

			c, err := client.New(cfg)

			if err != nil {
				log.Fatal(err)
			}

			kapi := client.NewKeysAPI(c)

			resp, err := kapi.Get(context.Background(), key, nil)

			if err != nil {
				log.Fatal(err)
			}

			mac_addr := string(resp.Node.Value)

			fmt.Printf("Mac address is: %s : %s \n", resp.Node.Key, mac_addr)

			fmt.Printf("pair Value %s \n", pair.Value)

			if err != nil {
				log.Fatal(err)
			}

			// after read the mac from etcd, start the container with the eth_peer mac address
			start_container(mac_addr)

			// run one time, stop here
			return
		}
	}
}
func start_container(mac_addr string) {
	// change the parameter if neccessry
	//cmd := "docker run -i -t -v /usr/local/var/run/openvswitch/vhost-user-1:/var/run/usvhost -v /dev/hugepages:/dev/hugepages dpdk-app-testpmd testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0," + mac_addr

	// run with the auto start
	//cmd := "docker run -i -t -v /usr/local/var/run/openvswitch/vhost-user-1:/var/run/usvhost -v /dev/hugepages:/dev/hugepages dpdk-app-testpmd testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -a --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0," + mac_addr

	// run with the auto start, txonly mode
	cmd := "docker run -i -t -v /usr/local/var/run/openvswitch/vhost-user-1:/var/run/usvhost -v /dev/hugepages:/dev/hugepages dpdk-app-testpmd testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -t --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0," + mac_addr

	exec_cmd(cmd)
}

func exec_cmd(cmd string) {
	log.Println("command is : ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("Error to exec CMD", cmd)
	}
	log.Println("Output of command:", string(out))
}
