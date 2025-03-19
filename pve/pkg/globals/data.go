package globals

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	CreateType  string = "create"
	DestroyType string = "destroy"
)
