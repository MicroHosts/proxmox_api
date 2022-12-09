package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "microhost_proxmox/internal/controller/http/v1"
	service2 "microhost_proxmox/internal/domain/service"
	"microhost_proxmox/internal/proxmox"
	"microhost_proxmox/pkg/mysql"
	"os"
)

func main() {
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	c, err := proxmox.NewClient(os.Getenv("PM_API_URL"), nil, os.Getenv("PM_HTTP_HEADERS"), tlsconf, "", 300)
	if err == nil {
		fmt.Print(err)
	}
	c.SetAPIToken("root@pam!test", "")
	db := mysql.NewClient("root", "example", "localhost", "3306", "test")
	service := service2.NewHostService(db, c)
	router := gin.Default()

	handler := v1.NewHostHandler(service)
	handler.Register(router)

	router.Run(":8080")
}
