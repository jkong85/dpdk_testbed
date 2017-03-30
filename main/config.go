package main

import (
	//"os"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	// We can register as many backends that are supported by libkv
	etcd.Register()

	server := "127.0.0.1:4001"

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{server},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store", kv)
	}

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:4001"},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	ip := "10.145.240.185"
	key := "/" + ip
	fileName := "mac.txt"

	// mac address read from the file
	mac_buf, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := kapi.Set(context.Background(), key, string(mac_buf), nil)

	if err != nil {
		log.Fatal(err)
	} else {
		//log.Printf("Set is done. Metadata is %q\n", resp)
		fmt.Println("Set is done. Metadata is %q\n", resp)
	}

}
