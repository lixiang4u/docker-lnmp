package model

type Label struct {
	Project    string `json:"project"`
	Service    string `json:"service"`
	Version    string `json:"version"`
	ConfigFile string `json:"config_file"`
	WorkingDir string `json:"working_dir"`
}
