package models

type CreateResp struct {
	VMID   string `json:"vm_id"`
	Email  string `json:"email"`
	IPAddr string `json:"ip_addr"`
}
