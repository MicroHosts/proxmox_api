package main

import (
	"crypto/tls"
	"log"
	"microhost_proxmox/proxmox"
	"os"
)

func main() {
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	c, err := proxmox.NewClient(os.Getenv("PM_API_URL"), nil, os.Getenv("PM_HTTP_HEADERS"), tlsconf, "", 300)
	c.SetAPIToken("", "")
	var jbody interface{}
	var vmr *proxmox.VmRef
	vmr = proxmox.NewVmRef(105)
	jbody, err = c.StopVm(vmr)
	failError(err)
	if jbody != nil {
		log.Println(jbody)
	}
}

func failError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
