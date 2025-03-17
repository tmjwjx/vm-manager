package repo

type CreateVM struct {
	VMID   string `json:"vm_id"`
	Email  string `json:"email"`
	IpAddr string `json:"ip_addr"`
}
