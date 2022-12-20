package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	config2 "microhost_proxmox/internal/config"
	service2 "microhost_proxmox/internal/domain/service"
	"microhost_proxmox/internal/proxmox"
	"microhost_proxmox/pkg/mysql"
	"os"
	"time"
)

var (
	Logger *log.Logger
)

func main() {
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	config := config2.GetConfig()
	c, err := proxmox.NewClient(config.Proxmox.PmApiUrl, nil, os.Getenv("PM_HTTP_HEADERS"), tlsconf, "", 300)
	if err != nil {
		fmt.Print(err)
		return
	}
	c.SetAPIToken(config.Proxmox.PmUser, config.Proxmox.PmToken)
	db := mysql.NewClient(config.MySQL.User, config.MySQL.Password, config.MySQL.Host,
		config.MySQL.Port, config.MySQL.DB)
	service := service2.NewHostService(db, c)
	ticker := time.NewTicker(time.Minute * 1)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			paymentCheck(db, service)
			deleteNoPayment(db, service)
			deletePaymentExpire(db)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func paymentCheck(db *sql.DB, service service2.HostService) {
	hosts, err := db.Query("SELECT id, hostId, userId FROM OrderHost WHERE rentDate < NOW()")
	if err != nil {
		return
	}
	for hosts.Next() {
		var id, hostId, userId string
		err = hosts.Scan(&id, &hostId, &userId)
		if err != nil {
			continue
		}
		add := time.Now().Add(time.Hour * 24 * 3)
		_, err = db.Exec("INSERT INTO NoPayOrderHost (hostId, userId, rentDate) VALUES (?, ?, ?)", hostId, userId, add)
		if err != nil {
			continue
		}
		_, err = db.Exec("DELETE FROM OrderHost WHERE id = ?", id)
		if err != nil {
			continue
		}
		host := service.GetHost(hostId)
		if host == nil {
			continue
		}
		err := service.Stop(host.Vimid)
		if err != nil {
			continue
		}
	}
}

func deleteNoPayment(db *sql.DB, service service2.HostService) {
	nopayhosts, err := db.Query("SELECT id, hostId, userId FROM NoPayOrderHost WHERE rentDate < NOW()")
	if err != nil {
		return
	}
	for nopayhosts.Next() {
		var id, hostId, userId string
		err = nopayhosts.Scan(&id, &hostId, &userId)
		_, err = db.Query(`UPDATE Host SET ready=$1 WHERE id = $2;`, false, hostId)
		_, err = db.Exec("DELETE FROM NoPayOrderHost WHERE id = ?", id)
		if err != nil {
			continue
		}
	}
}

func deletePaymentExpire(db *sql.DB) {
	db.Exec("DELETE FROM `Payment` WHERE `createdAt` <= ( CURDATE() - INTERVAL 2 DAY ) AND `paid`=false;")
}
