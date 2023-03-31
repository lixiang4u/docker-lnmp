package model

type Container struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Image     string   `json:"image"`
	ImageId   string   `json:"image_id"`
	Labels    Label    `json:"labels"`
	State     string   `json:"state"`
	Ports     []string `json:"ports"`
	Status    string   `json:"status"`
	CreatedAt int64    `json:"created_at"`
}
