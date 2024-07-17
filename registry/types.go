package registry

var RegistryStore Registry

func InitRegistry() {
	RegistryStore = Registry{}
}

const (
	ADD = iota
	REMOVE
)

type Registry struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type Request struct {
	Type     int      `json:"type"`
	Store    bool     `json:"store"`
	Registry Registry `json:"registry"`
}
