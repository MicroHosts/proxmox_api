package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "microhost_proxmox/internal/controller/http/v1"
	service2 "microhost_proxmox/internal/domain/service"
	"microhost_proxmox/internal/proxmox"
	"microhost_proxmox/pkg/mysql"
	"os"
	"time"
)

func main() {
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	c, err := proxmox.NewClient(os.Getenv("PM_API_URL"), nil, os.Getenv("PM_HTTP_HEADERS"), tlsconf, "", 300)
	if err == nil {
		fmt.Print(err)
		return
	}
	c.SetAPIToken(os.Getenv("PM_USER"), os.Getenv("PM_TOKEN"))
	db := mysql.NewClient(os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_IP"),
		os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
	service := service2.NewHostService(db, c)
	router := gin.Default()

	handler := v1.NewHostHandler(service)
	handler.Register(router)
	ticker := time.NewTicker(time.Hour * 5)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				paymentCheck(c, db, service)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	router.Run(":8080")
}

func paymentCheck(client *proxmox.Client, db *sql.DB, service service2.HostService) {
	hosts, err := db.Query("SELECT id, hostId, userId FROM OrderHost WHERE rentDate < NOW()")
	if err != nil {
		return
	}
	for hosts.Next() {
		var id, hostId, userId string
		err = hosts.Scan(&id, &hostId, &userId)
		if err != nil {
			return
		}
		add := time.Now().Add(time.Hour * 24 * 30)
		_, err = db.Exec("INSERT INTO NoPayOrderHost (hostId, userId, rentDate) VALUES (?, ?, ?)", hostId, userId, add)
		if err != nil {
			return
		}
		_, err = db.Exec("DELETE FROM OrderHost WHERE id = ?", id)
		if err != nil {
			return
		}
		host := service.GetHost(hostId)
		if host == nil {
			return
		}
		//convert host.Vimid to int
		service.Stop(host.Vimid)
	}
}
