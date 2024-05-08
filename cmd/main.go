package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/luthermonson/go-proxmox"
)

func main() {
	proxmoxAddr := "http://" + os.Getenv("PROXMOX_ADDRESS") + ":8006/api2/json"

	insecureHTTPClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	tokenID := os.Getenv("PROXMOX_TOKEN_ID")
	token := os.Getenv("PROXMOX_TOKEN_SECRET")
	client := proxmox.NewClient(proxmoxAddr,
		proxmox.WithHTTPClient(&insecureHTTPClient),
		proxmox.WithAPIToken(tokenID, token),
	)

	// create vm

	// 1. get node
	node, err := client.Node(context.TODO(), "pve1")
	if err != nil {
		fmt.Println("Node Error: ", err)
		return
	}

	// 2. create vm
	vmID, err := strconv.Atoi(os.Getenv("VM_ID"))
	if err != nil {
		fmt.Println("VM ID Error: ", err)
		return
	}

	task, err := node.NewVirtualMachine(
		context.TODO(),
		vmID,
		proxmox.VirtualMachineOption{
			Name:  "name",
			Value: os.Getenv("VM_NAME"),
		},
		proxmox.VirtualMachineOption{
			Name:  "cores",
			Value: os.Getenv("VM_CORES"),
		},
		proxmox.VirtualMachineOption{
			Name:  "memory",
			Value: os.Getenv("VM_MEMORY"),
		},
	)

	if err != nil {
		fmt.Println("Virtual Machine Error: ", err)
		return
	}

	// 3. wait for task to complete
	err = task.Wait(context.TODO(), 5, 10)
	if err != nil {
		fmt.Println("Task Error: ", err)
		return
	}

}
