package service

import (
	"database/sql"
	"fmt"
	"microhost_proxmox/internal/domain/entity"
	"microhost_proxmox/internal/proxmox"
)

type HostService interface {
	GetHost(id string) *entity.Host
	Stop(vimId int) error
	Start(vimId int) error
	Reboot(vimId int) error
	Shutdown(vimId int) error
}

type hostService struct {
	client  *sql.DB
	proxmox *proxmox.Client
}

func NewHostService(client *sql.DB, proxmox *proxmox.Client) HostService {
	return &hostService{
		client:  client,
		proxmox: proxmox,
	}
}

func (h hostService) GetHost(id string) *entity.Host {
	res, err := h.client.Query("SELECT id, name, vimid, price FROM Host WHERE id= ?", id)

	if err != nil {
		//TODO log error
		return nil
	}

	if res.Next() {
		var host entity.Host
		err := res.Scan(&host.Id, &host.Name, &host.Vimid, &host.Price)
		if err != nil {
			fmt.Print(err)
			return nil
		}
		return &host
	} else {
		return nil
	}
}

func (h hostService) Stop(vimId int) error {
	vmr := proxmox.NewVmRef(vimId)
	_, err := h.proxmox.StopVm(vmr)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) Start(vimId int) error {
	vmr := proxmox.NewVmRef(vimId)
	_, err := h.proxmox.StartVm(vmr)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) Reboot(vimId int) error {
	vmr := proxmox.NewVmRef(vimId)
	_, err := h.proxmox.ResetVm(vmr)
	if err != nil {
		return err
	}
	return nil
}

func (h hostService) Shutdown(vimId int) error {
	vmr := proxmox.NewVmRef(vimId)
	_, err := h.proxmox.ShutdownVm(vmr)
	if err != nil {
		return err
	}
	return nil
}
