package model

import (
	"fmt"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/util"
	"github.com/spf13/viper"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type VirtualHost struct {
	Id      string `yaml:"id" json:"id" toml:"id"`                   // ID全局唯一
	Name    string `yaml:"name" json:"name" toml:"name"`             // 不带特殊字符的名称，全局唯一
	Domain  string `yaml:"domain" json:"domain" toml:"domain"`       // 虚拟主机域名
	Root    string `yaml:"root" json:"root" toml:"root"`             // 项目根目录（本地机器目录），用于docker的volumes映射
	WebRoot string `yaml:"web_root" json:"web_root" toml:"web_root"` // 项目web服务的根目录（本地机器目录），相对于root位置，一般为："/", "/public"这类
	Port    int    `yaml:"port" json:"port" toml:"port"`             // 虚拟主机端口
}

func InitVirtualHostConfig() error {
	// 需要在配置docker前设置默认虚拟主机
	file, err := os.OpenFile("config.json", os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	viper.SetConfigFile("config.json")
	err = viper.ReadInConfig()
	if err != nil {
		_ = file.Truncate(0)
		_, _ = file.WriteString("{}")
		log.Fatalln("[viper.ReadConfig error]", err.Error())
		return err
	}
	var hosts = []VirtualHost{
		{
			Id:      util.StringHash("default.me"),
			Name:    "default",
			Domain:  "default.me",
			Root:    filepath.Join(util.AppDirectory(), "dockerfile/nginx/html"),
			WebRoot: "",
			Port:    0,
		},
	}

	viper.SetDefault("_READE_ME", "该配置自动生成，请勿修改")
	viper.SetDefault("app", config.AppName)
	viper.SetDefault("hosts", hosts)
	err = viper.WriteConfig()
	if err != nil {
		log.Fatalln("[viper.WriteConfig error]", err.Error())
		return err
	}
	return nil
}

func ConvertToVirtualHost(v any) []VirtualHost {
	if v, ok := v.([]VirtualHost); ok {
		return v
	}
	if tmpHosts, ok := v.([]interface{}); ok {
		var tmpList []VirtualHost
		for _, tmpHost := range tmpHosts {
			if h, ok := tmpHost.(map[string]interface{}); ok {
				tmpList = append(tmpList, VirtualHost{
					Id:      h["id"].(string),
					Name:    h["name"].(string),
					Domain:  h["domain"].(string),
					Root:    h["root"].(string),
					WebRoot: h["web_root"].(string),
					Port:    0,
				})
			}
		}
		return tmpList
	}
	return []VirtualHost{}
}

func UpdateVirtualHost(hosts []VirtualHost) DockerComposeTpl {
	var dc = InitComposeConfig()
	for _, host := range hosts {
		if _, ok := dc.Services["nginx"]; ok {
			var nginxService = dc.Services["nginx"].(DockerComposeServiceTpl)
			var pWWW = fmt.Sprintf(
				"%s:%s",
				host.Root,
				fmt.Sprintf("/apps/www/%s:ro,bind", host.Domain),
			)
			nginxService.Volumes = append(nginxService.Volumes, pWWW)
			// 将修改后的 nginxService 再赋回原来的 map 中
			dc.Services["nginx"] = nginxService
		}
		if _, ok := dc.Services["php"]; ok {
			var phpService = dc.Services["php"].(DockerComposeServiceTpl)
			var pWWW = fmt.Sprintf(
				"%s:%s",
				host.Root,
				fmt.Sprintf("/apps/www/%s:ro,bind", host.Domain),
			)
			phpService.Volumes = append(phpService.Volumes, pWWW)
			// 将修改后的 nginxService 再赋回原来的 map 中
			dc.Services["php"] = phpService
		}

		// 生成nginx配置文件
		updateNginxVirtualHostConfig(host)
	}

	return dc
}

func updateNginxVirtualHostConfig(host VirtualHost) {
	p, err := parser.NewParser(filepath.Join(util.AppDirectory(), "dockerfile/nginx/config/default.tpl"))
	if err != nil {
		log.Fatalln("[parser.NewParser error]", err.Error())
		return
	}
	c := p.Parse()

	var directives []gonginx.IDirective
	directives = c.FindDirectives("server_name")
	for _, directive := range directives {
		directive.GetParameters()[0] = host.Domain
	}
	directives = c.FindDirectives("root")
	for _, directive := range directives {
		directive.GetParameters()[0] = filepath.Join(
			fmt.Sprintf("/apps/www/%s", host.Domain),
			host.WebRoot,
		)
		directive.GetParameters()[0] = filepath.ToSlash(directive.GetParameters()[0])
	}
	// golang中使用filepath 生成linux中的目录路径
	var pConfig = filepath.Join(
		util.AppDirectory(),
		"dockerfile/nginx/config",
		fmt.Sprintf("%s.conf", host.Domain),
	)
	file, err := os.OpenFile(pConfig, os.O_CREATE|os.O_RDWR|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		log.Fatalln("[os.Open error]", err.Error())
		return
	}
	_, err = file.WriteString(gonginx.DumpBlock(c.Block, gonginx.IndentedStyle))
	if err != nil {
		log.Fatalln("[os.WriteString error]", err.Error())
		return
	}
}
