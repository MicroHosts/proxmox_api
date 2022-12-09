package v1

import (
	"github.com/gin-gonic/gin"
	"microhost_proxmox/internal/api"
	"microhost_proxmox/internal/controller/dto"
	"microhost_proxmox/internal/domain/service"
	"net/http"
	"strconv"
)

type handler struct {
	service service.HostService
}

func NewHostHandler(service service.HostService) api.Handler {
	return &handler{service: service}
}

func (h handler) Register(router *gin.Engine) {
	router.GET("/v1/host", h.Get)
	router.GET("/v1/host/stop", h.Stop)
	router.GET("/v1/host/start", h.Start)
	router.GET("/v1/host/reboot", h.Reboot)
	router.GET("/v1/host/shutdown", h.Shutdown)
}

func (h *handler) Get(c *gin.Context) {
	var json dto.HostGetDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	host := h.service.GetHost(json.ID)
	c.JSON(http.StatusOK, host)
}

func (h *handler) Stop(c *gin.Context) {
	var json dto.VMIDDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vimId, err := strconv.Atoi(json.Vmid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите VIMID"})
		return
	}
	h.service.Stop(vimId)
	c.JSON(http.StatusOK, gin.H{"message": "Сервер остановлен"})
}

func (h *handler) Start(c *gin.Context) {
	var json dto.VMIDDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vimId, err := strconv.Atoi(json.Vmid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите VIMID"})
		return
	}
	h.service.Start(vimId)
	c.JSON(http.StatusOK, gin.H{"message": "Сервер включен"})
}

func (h *handler) Reboot(c *gin.Context) {
	var json dto.VMIDDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vimId, err := strconv.Atoi(json.Vmid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите VIMID"})
		return
	}
	h.service.Reboot(vimId)
	c.JSON(http.StatusOK, gin.H{"message": "Сервер перезапущен"})
}

func (h *handler) Shutdown(c *gin.Context) {
	var json dto.VMIDDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vimId, err := strconv.Atoi(json.Vmid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите VIMID"})
		return
	}
	h.service.Shutdown(vimId)
	c.JSON(http.StatusOK, gin.H{"message": "Сервер выключен"})
}
