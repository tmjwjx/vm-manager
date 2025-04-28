package models

type CreateReq struct {
	Email string `json:"email"`
}

type StartReq struct {
	VMID string `json:"vm_id"`
}

type RebootReq struct {
	VMID string `json:"vm_id"`
}

type GetVMInfoReq struct {
	VMID string `json:"vm_id"`
}
