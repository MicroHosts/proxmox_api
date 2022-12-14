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
	//go func() {
	//
	//}()
	for {
		select {
		case <-ticker.C:
			paymentCheck(db, service)
			deleteNoPayment(db, service)
		case <-quit:
			ticker.Stop()
			return
		}
	}
	//Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//websocketproxy.SetLogger(Logger)
	//websocketproxy.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	//wp, err := websocketproxy.NewProxy("wss://192.168.0.2:8006/api2/json/nodes/dc1/qemu/112/vncwebsocket?port=5900&"+
	//	"vncticket=PVEVNC:639941AE::uJJnBqe5pDiBqWdYyXdvxgY6PDJwwzfd3LOFVQbpQLAH0jFujPfiEYYXcYu9zTXoUbbb96q4oeCjU4So6M6vADSoYJ7BAMfo1sARk8hTRoLmWH2TIYWDGovV5uAiMEYWqqSfqBG/5QeCsE+swouvtHRZG+KVHSIa1ELm3fbG3TChscKZyxpLE+vZI+l5S26hZS6Sau+su0tCc0nbE1ZKKLoMz7bFlWc3nzFyUuCNYUVRUfodo3EgLs2M5o1O7bDceRXwFX8p8Nro21HWAC5SCvf7JGQZwUzahMbEXUUXtwQ/KR/l4MXvi4O2iHhLlyh0WVhuXm1ZtfroaA2XMiAj5Q==", func(r *http.Request) error {
	//	// Permission to verify
	//	cookie := http.Cookie{Name: "PVEAuthCookie", Value: "PVE%3Aroot@pam%3A63994128%3A%3AMKVz2uh4f3bEIaZmZp+p6mlglIe6PtqHtcHs1RxFSG/G86Wz2qc5ox37LzE7fIYleBa+VphbIEkG/s95nQw5mbMxWDASkraYyGKZDWeFr/gqvh3yn8HKW62RT4QOTp6RWt4qHsMmg4QHSoDsTBxoBv30jbh/w1fX/OAGOE/qj80pMUDEFKokgom0rPIblFRXdUNaCDBjjpEdtpRKywkJ0MliqKAVNzdfDx+S8jtUAOSCNaevCA6AXdxvlbMDhilkTf9BcPkyoIiKQ16MglwbkNzfR2R/gINV9wsJrO3vL21fFuRf2hco6barzsFpSoT5niIUtxIfV1j/I7/AbUn/XA%3D%3D"}
	//	r.AddCookie(&cookie)
	//	// Source of disguise
	//	r.Header.Set("Origin", "https://192.168.0.2:8006")
	//	return nil
	//})

	//if err != nil {
	//	fmt.Print("asdsad")
	//	return
	//}
	//
	//http.HandleFunc("/wsproxy", wp.Proxy)
	//http.ListenAndServe(":9696", nil)
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
