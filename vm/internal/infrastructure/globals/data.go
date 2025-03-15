package globals

type Data struct {
	Type DataType `json:"type"`
	Data []byte   `json:"data"`
}

type DataType string

const (
	CreateType  DataType = "create"
	DestroyType DataType = "destroy"
)
