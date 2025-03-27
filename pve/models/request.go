package models

type CreateReq struct {
	Email string `json:"email"`
}

type StartReq struct {
	VMID string `json:"vm_id"`
}
