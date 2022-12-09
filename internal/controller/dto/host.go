package dto

type HostGetDTO struct {
	ID string `json:"id" binding:"required"`
}

type VMIDDTO struct {
	Vmid string `json:"vmid" binding:"required"`
}
