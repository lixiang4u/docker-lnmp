package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/model"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type HostController struct {
	//ctx.JSON(http.StatusOK,gin.H{"code":200,"msg":"","data":nil})
}

func (x HostController) Init(ctx *gin.Context) {

	var hosts = []model.VirtualHost{
		{
			Name:    "api",
			Domain:  "api.local.me",
			Root:    "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\dockerfile\\nginx\\html",
			WebRoot: "",
			Port:    0,
		},
		{
			Name:    "mobile",
			Domain:  "m.local.me",
			Root:    "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\dockerfile\\nginx\\html",
			WebRoot: "/mm",
			Port:    0,
		},
		{
			Name:    "download",
			Domain:  "d.local.me",
			Root:    "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\dockerfile\\nginx\\html",
			WebRoot: "",
			Port:    0,
		},
		{
			Name:    "tv",
			Domain:  "tv.local.me",
			Root:    "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\dockerfile\\nginx\\html",
			WebRoot: "",
			Port:    0,
		},
	}

	var dc = x.updateVirtualHost(hosts)

	out, err := yaml.Marshal(dc)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": err.Error(), "data": nil})
		return
	}

	ctx.String(http.StatusOK, string(out))

	return
	//var dc = model.DockerComposeTpl{
	//	Version:  "3",
	//	Networks: map[string]interface{}{"default-network": nil},
	//	Volumes: map[string]interface{}{
	//		"default_mariadb_data": nil,
	//		"default_redis_data":   nil,
	//	},
	//	Services: map[string]interface{}{
	//		"nginx": model.DockerComposeServiceTpl{
	//			ContainerName: "lnmp-nginx",
	//			Image:         "nginx:stable-alpine",
	//			Networks:      []string{"default-network"},
	//			DependsOn:     []string{"php", "mariadb", "redis"},
	//			Ports:         []string{"80:80", "443:443"},
	//			Volumes: []string{
	//				"D:/ProgramData/docker-lamp/nginx/log:/var/log/nginx",
	//			},
	//		},
	//		"php": model.DockerComposeServiceTpl{
	//			ContainerName: "lnmp-php72",
	//			Build:         map[string]interface{}{"dockerfile": "files/dockerfile/lamp-php-fpm"},
	//			Networks:      []string{"default-network"},
	//			Volumes: []string{
	//				"D:/repo/aidun/74cms:/apps/www/rencai.local.me",
	//			},
	//		},
	//		"mariadb": model.DockerComposeServiceTpl{
	//			ContainerName: "lnmp-mariadb",
	//			Image:         "mariadb:10.7.8",
	//			Networks:      []string{"default-network"},
	//			Environment:   []string{"MARIADB_ROOT_PASSWORD=123456"},
	//			Volumes: []string{
	//				"default_mariadb_data:/var/lib/mysql",
	//			},
	//		},
	//		"redis": model.DockerComposeServiceTpl{
	//			ContainerName: "lnmp-redis",
	//			Image:         "redis:latest",
	//			Networks:      []string{"default-network"},
	//			Volumes: []string{
	//				"default_redis_data:/data",
	//			},
	//		},
	//		"phpmyadmin": model.DockerComposeServiceTpl{
	//			ContainerName: "lnmp-sql-admin",
	//			Image:         "phpmyadmin:latest",
	//			Networks:      []string{"default-network"},
	//			Environment: []string{
	//				"PMA_ARBITRARY=1",
	//				"PMA_HOST=lnmp-mariadb",
	//				"PMA_PORT=3306",
	//				"UPLOAD_LIMIT=512M",
	//			},
	//		},
	//	},
	//}
	//
	////if v, ok := dc.Networks["default-network"]; !ok {
	////	log.Printf("！dc.Networks[\"ss\"]")
	////} else {
	////	log.Printf("vvv: ", v)
	////}
	//
	//ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": nil, "data": dc})
}

func (x HostController) List(ctx *gin.Context) {
	file, err := os.Open("D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\docker-compose.yaml")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	bs, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	m := make(map[string]interface{})
	err = yaml.Unmarshal(bs, &m)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	log.Printf("[=====>1 ]%#v", m)

	bs, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	log.Printf("[=====>2 ] %s", string(bs))

	//ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": nil, "data": string(bs)})
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": nil, "data": m})

	return
}
func (x HostController) Show(context *gin.Context) {
	//
}
func (x HostController) Update(context *gin.Context) {
	//
}
func (x HostController) Create(context *gin.Context) {
	//
}
func (x HostController) Delete(context *gin.Context) {
	//
}

func (x HostController) parseYAML() error {
	//file, err := os.Open("")
	//if err != nil {
	//	return err
	//}
	//decoder := yaml.NewDecoder(file)
	//decoder.

	return nil
}

func (x HostController) initConfig() model.DockerComposeTpl {
	var dc = model.DockerComposeTpl{
		Version:  "3",
		Networks: map[string]interface{}{"default-network": nil},
		Volumes: map[string]interface{}{
			"default_mariadb_data": nil,
			"default_redis_data":   nil,
		},
		Services: map[string]interface{}{
			"nginx": model.DockerComposeServiceTpl{
				ContainerName: "lnmp-nginx",
				Image:         "nginx:stable-alpine",
				Networks:      []string{"default-network"},
				DependsOn:     []string{"php", "mariadb", "redis"},
				Ports:         []string{"80:80", "443:443"},
				Volumes:       []string{},
			},
			"php": model.DockerComposeServiceTpl{
				ContainerName: "lnmp-php72",
				Build:         map[string]interface{}{"dockerfile": "dockerfile/script/lamp-php-fpm"},
				Networks:      []string{"default-network"},
				Volumes:       []string{},
			},
			"mariadb": model.DockerComposeServiceTpl{
				ContainerName: "lnmp-mariadb",
				Image:         "mariadb:10.7.8",
				Networks:      []string{"default-network"},
				Environment:   []string{"MARIADB_ROOT_PASSWORD=123456"},
				Volumes: []string{
					"default_mariadb_data:/var/lib/mysql",
				},
			},
			"redis": model.DockerComposeServiceTpl{
				ContainerName: "lnmp-redis",
				Image:         "redis:latest",
				Networks:      []string{"default-network"},
				Volumes: []string{
					"default_redis_data:/data",
				},
			},
			"phpmyadmin": model.DockerComposeServiceTpl{
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
		var nginxService = dc.Services["nginx"].(model.DockerComposeServiceTpl)
		var pLog = fmt.Sprintf("%s:%s", filepath.Join(x.currentDirectory(), "dockerfile/nginx/log"), "/var/log/nginx")
		var pConfig = fmt.Sprintf("%s:%s", filepath.Join(x.currentDirectory(), "dockerfile/nginx/config"), "/etc/nginx/conf.d")
		var pWWW = fmt.Sprintf("%s:%s", filepath.Join(x.currentDirectory(), "dockerfile/nginx/html"), "/apps/www/default.me:ro,bind")
		nginxService.Volumes = append(nginxService.Volumes, pLog)
		nginxService.Volumes = append(nginxService.Volumes, pConfig)
		nginxService.Volumes = append(nginxService.Volumes, pWWW)
		// 将修改后的 nginxService 再赋回原来的 map 中
		dc.Services["nginx"] = nginxService
	}

	return dc
}

func (x HostController) updateVirtualHost(hosts []model.VirtualHost) model.DockerComposeTpl {
	var dc = x.initConfig()
	for _, host := range hosts {
		if _, ok := dc.Services["nginx"]; ok {
			var nginxService = dc.Services["nginx"].(model.DockerComposeServiceTpl)
			var pWWW = fmt.Sprintf(
				"%s:%s",
				filepath.Join(x.currentDirectory(), "dockerfile/nginx/html"),
				fmt.Sprintf("/apps/www/%s:ro,bind", host.Domain),
			)
			nginxService.Volumes = append(nginxService.Volumes, pWWW)
			// 将修改后的 nginxService 再赋回原来的 map 中
			dc.Services["nginx"] = nginxService
		}
		if _, ok := dc.Services["php"]; ok {
			var phpService = dc.Services["php"].(model.DockerComposeServiceTpl)
			var pWWW = fmt.Sprintf(
				"%s:%s",
				filepath.Join(x.currentDirectory(), "dockerfile/nginx/html"),
				fmt.Sprintf("/apps/www/%s:ro,bind", host.Domain),
			)
			phpService.Volumes = append(phpService.Volumes, pWWW)
			// 将修改后的 nginxService 再赋回原来的 map 中
			dc.Services["php"] = phpService
		}

		// 生成nginx配置文件
		x.updateNginxHostConfig(host)
	}

	return dc
}

func (x HostController) updateNginxHostConfig(host model.VirtualHost) {
	p, err := parser.NewParser(filepath.Join(x.currentDirectory(), "dockerfile/nginx/config/default.tpl"))
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
		x.currentDirectory(),
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

//func (x HostController) initNginxConfig() string {
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		return ""
//	}
//	p, err := parser.NewParser(filepath.Join(x.currentDirectory(), "/../../dockerfile/nginx/config/default.me.tpl"))
//	if err != nil {
//		log.Printf("[err] %s", err.Error())
//		return ""
//	}
//	c := p.Parse()
//	for _, directive := range c.FindDirectives("root") {
//		directive.GetParameters()[0] = "/path/to/some/where"
//	}
//	for _, directive := range c.FindDirectives("listen") {
//		directive.GetParameters()[0] = "8070"
//	}
//	for _, directive := range c.FindDirectives("server_name") {
//		directive.GetParameters()[0] = "debug.local.me"
//	}
//
//	log.Println()
//	log.Println()
//	log.Println("======================================================")
//
//	for _, directive := range c.FindDirectives("location") {
//		log.Printf("[GetName]  %#v", directive.GetName())
//		log.Printf("[GetParameters] %#v", directive.GetParameters())
//		log.Printf("[GetComment] %#v", x.toJson(directive.GetComment()))
//		log.Printf("[GetBlock] %#v", directive.GetBlock())
//		for _, d2 := range directive.GetBlock().GetDirectives() {
//			log.Printf("      [BLOCK.NAME] %#v", d2.GetName())
//			log.Printf("      [BLOCK.PARA] %#v", d2.GetParameters())
//			log.Printf("      [BLOCK.COMM] %#v", d2.GetComment())
//			log.Printf("      [BLOCK.BLOC] %#v", d2.GetBlock())
//		}
//
//		log.Println()
//		log.Println()
//		log.Println("===================")
//	}
//
//	// fmt.Println(gonginx.DumpBlock(c.Block, gonginx.IndentedStyle))
//
//	return dir
//}

func (x HostController) currentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return filepath.Dir(filepath.Dir(dir))
}

func (x HostController) toJson(v any) string {
	bs, _ := json.MarshalIndent(v, "", "\t")
	log.Println(string(bs))
	return string(bs)
}
