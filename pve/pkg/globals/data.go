package globals

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	CreateType    string = "create"
	DestroyType   string = "destroy"
	StartType     string = "start"
	RebootType    string = "reboot"
	GetVMInfoType string = "get_vm_info"
)
