package model

type DockerComposeTpl struct {
	Version  string                 `yaml:"version" json:"version"`
	Services map[string]interface{} `yaml:"services" json:"services"`
	Networks map[string]interface{} `yaml:"networks" json:"networks"`
	Volumes  map[string]interface{} `yaml:"volumes" json:"volumes"`
}

type DockerComposeServiceTpl struct {
	ContainerName string                 `yaml:"container_name" json:"container_name"`
	Build         map[string]interface{} `yaml:"build" json:"build"`
	Image         string                 `yaml:"image" json:"image"`
	Networks      []string               `yaml:"networks" json:"networks"`
	Volumes       []string               `yaml:"volumes" json:"volumes"`
	DependsOn     []string               `yaml:"depends_on" json:"depends_on"`
	Ports         []string               `yaml:"ports" json:"ports"`
	Environment   []string               `yaml:"environment" json:"environment"`
}
