package virtualMachine

/*
接受pve项目的响应的结构体
*/

const (
	CreateType  string = "create"
	DestroyType string = "destroy"
	StartType   string = "start"
	StopType    string = "stop"
	RenewType   string = "renew"
	GetInfoType string = "getInfo"
)

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type CreateResp struct {
	VMID  string `json:"vmid"`
	Email string `json:"email"`
}
