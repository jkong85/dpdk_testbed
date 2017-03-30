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
	// the main structure for the DPDK demo
	/*
			Need to add dpdkport ( add here or use the CNI ??)
			kubectl dpdk-1.yaml
			For pod 1: start testpmd with the -a option, and it will send the mac to the file
				Then on vm 1: write the mac address to etcd server
		    then kubectl dpdk-2.yaml
			get the pod 1's mac address by reading the etcd server
			For pod 2: start testpmd with the -t -eth_address option. it will auto start traffic to the pod 1.
	*/

	// Is there need to add DPDK port here or in the CNI??

	// start the pod-1
	yamlfile_pod_rx := "dpdk-pod-rx.yaml"
	yamlfile_pod_tx := "dpdk-pod-tx.yaml"

	createPod(yamlfile_pod_rx)
	createPod(yamlfile_pod_tx)

	pod_name_rx, pod_ip_rx, pod_hip_rx := getPodInfo("dpdk-pod-rx")
	pod_name_tx, pod_ip_tx, pod_hip_tx := getPodInfo("dpdk-pod-tx")
	log.Println(" Pod RX info:" + pod_name_rx + " " + pod_ip_rx + " " + pod_hip_rx)
	log.Println(" Pod RX info:" + pod_name_tx + " " + pod_ip_tx + " " + pod_hip_tx)

	if pod_hip_rx == pod_hip_tx {
		log.Println("Do NOT support the test on the same node!")
		return
	}

	startTestPMDRX(pod_name_rx)

	// wait several seconds for writing the mac address to the etcd server
	time.Sleep(3 * time.Second)

	key := "/" + string(pod_hip_rx)
	log.Println("etcd key is: " + key)

	macAddress := getMacAddress(key)
	log.Println("dst mac address is: ", macAddress)

	// then we start the testPMD on pod 2
	startTestPMDTX(pod_name_tx, macAddress)

	// check the result here
	log.Println("Next is to check the traffic statistics!")
}

// get pod's info (name, ip, hip)
func getPodInfo(podName string) (string, string, string) {

	// getPodHostIP.sh
	// 1: get pod name, 2: get pod's ip, 3: get host IP
	cmd := "./getPodHostIP.sh " + podName + " 1 "
	name := exec_cmd(cmd)

	cmd = "./getPodHostIP.sh " + podName + " 2 "
	ip := exec_cmd(cmd)

	cmd = "./getPodHostIP.sh " + podName + " 3 "
	hip := exec_cmd(cmd)

	return name, ip, hip
}

// return the pod's pid
func createPod(yamlfile string) {

	delcmd := "kubectl delete -f " + yamlfile
	exec_cmd(delcmd)
	time.Sleep(10 * time.Second)

	createcmd := "kubectl create -f " + yamlfile
	exec_cmd(createcmd)

	time.Sleep(10 * time.Second)

	if check_pods_created() {
		log.Println("pod creating fails!")
	}

	getcmd := "kubectl get pods -o wide"
	result := exec_cmd(getcmd)
	log.Println("create the pod-1: %s", result)

}

// set the RX pod
func startTestPMDRX(podID string) {
	pmdcmd := " testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -a --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac"
	cmd := "kubectl exec " + podID + pmdcmd
	exec_cmd(cmd)
}

// set the TX pod
/*
run with the auto start, txonly mode
cmd := "docker run -i -t -v /usr/local/var/run/openvswitch/vhost-user-1:/var/run/usvhost -v /dev/hugepages:/dev/hugepages dpdk-app-testpmd testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -t --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0," + mac_addr
*/
func startTestPMDTX(podID string, mac_addr string) {
	pmdcmd := " testpmd -l 6-7 -n 4 -m 1024 --no-pci --vdev=virtio_user0,path=/var/run/usvhost -- -i -t --txqflags=0xf00 --disable-hw-vlan --forward-mode=mac --eth-peer=0," + mac_addr
	cmd := "kubectl exec " + podID + pmdcmd
	exec_cmd(cmd)
}

func check_pods_created() bool {
	return true
}

// get the mac address by reading the etcd server
func getMacAddress(key string) string {

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
			// run one time, stop here
			return mac_addr
		}
	}
}

func exec_cmd(cmd string) string {
	log.Println("command is : ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("Error to exec CMD", cmd)
	}
	log.Println("Output of command:", string(out))

	return string(out)
}
