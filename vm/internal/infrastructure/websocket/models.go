package websocket

/*
对pve项目发送的请求的结构体
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

type CreateVMReq struct {
	Email string `json:"email"`
}
type CreateResp struct {
	VMID  string `json:"vmid"`
	Email string `json:"email"`
}
