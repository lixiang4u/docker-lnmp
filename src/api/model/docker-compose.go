package model

import (
	"fmt"
	"github.com/lixiang4u/docker-lnmp/util"
	"path/filepath"
)

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

func InitComposeConfig() DockerComposeTpl {
	var dc = DockerComposeTpl{
		Version:  "3",
		Networks: map[string]interface{}{"default-network": nil},
		Volumes: map[string]interface{}{
			"default_mariadb_data": nil,
			"default_redis_data":   nil,
		},
		Services: map[string]interface{}{
			"nginx": DockerComposeServiceTpl{
				ContainerName: "lnmp-nginx",
				Image:         "nginx:stable-alpine",
				Networks:      []string{"default-network"},
				DependsOn:     []string{"php", "mariadb", "redis"},
				Ports:         []string{"80:80", "443:443"},
				Volumes:       []string{},
			},
			"php": DockerComposeServiceTpl{
				ContainerName: "lnmp-php72",
				Build:         map[string]interface{}{"dockerfile": "common/script/lamp-php-fpm"},
				Image:         "php-fpm:72",
				Networks:      []string{"default-network"},
				Volumes:       []string{},
			},
			"mariadb": DockerComposeServiceTpl{
				ContainerName: "lnmp-mariadb",
				Image:         "mariadb:10.7.8",
				Networks:      []string{"default-network"},
				Environment:   []string{"MARIADB_ROOT_PASSWORD=123456"},
				Volumes: []string{
					"default_mariadb_data:/var/lib/mysql",
				},
			},
			"redis": DockerComposeServiceTpl{
				ContainerName: "lnmp-redis",
				Image:         "redis:latest",
				Networks:      []string{"default-network"},
				Volumes: []string{
					"default_redis_data:/data",
				},
			},
			"phpmyadmin": DockerComposeServiceTpl{
				ContainerName: "lnmp-sql-admin",
				Image:         "phpmyadmin:latest",
				Networks:      []string{"default-network"},
				Environment: []string{
					"PMA_ARBITRARY=1",
					"PMA_HOST=lnmp-mariadb",
					"PMA_PORT=3306",
					"UPLOAD_LIMIT=512M",
				},
			},
		},
	}

	if _, ok := dc.Services["nginx"]; ok {
		var nginxService = dc.Services["nginx"].(DockerComposeServiceTpl)
		var pLog = fmt.Sprintf("%s:%s", filepath.Join(util.AppDirectory(), "common/nginx/log"), "/var/log/nginx")
		var pConfig = fmt.Sprintf("%s:%s", filepath.Join(util.AppDirectory(), "common/nginx/config"), "/etc/nginx/conf.d")
		nginxService.Volumes = append(nginxService.Volumes, pLog)
		nginxService.Volumes = append(nginxService.Volumes, pConfig)
		// 将修改后的 nginxService 再赋回原来的 map 中
		dc.Services["nginx"] = nginxService
	}

	return dc
}

func ComposeDown(configFile string) ([]byte, error) {
	return util.Exec(fmt.Sprintf(
		"%s && %s",
		fmt.Sprintf("cd %s", filepath.Dir(configFile)),
		"docker compose down",
	))
}

func ComposeUp(configFile string) ([]byte, error) {
	return util.Exec(fmt.Sprintf(
		"%s && %s",
		fmt.Sprintf("cd %s", filepath.Dir(configFile)),
		fmt.Sprintf("docker compose -f %s up --force-recreate -d", configFile),
	))
}

func ComposeDownUp(configFile string) ([]byte, error) {
	return util.Exec(fmt.Sprintf(
		"%s && %s && %s",
		fmt.Sprintf("cd %s", filepath.Dir(configFile)),
		"docker compose down",
		fmt.Sprintf("docker compose -f %s up --force-recreate -d", configFile),
	))
}
