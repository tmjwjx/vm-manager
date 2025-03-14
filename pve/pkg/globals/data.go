package globals

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	CreateType = "create"
	DeleteType = "delete"
)
