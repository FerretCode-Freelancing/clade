package containers

const (
	CREATE = iota
	DELETE
	UPDATE
)

type Request struct {
	Type     int    `json:"type"`
	Name     string `json:"name"`
	Image    string `json:"image_name"`
	ImageURL string `json:"image"`
}
