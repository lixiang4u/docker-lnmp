package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/model"
	"github.com/spf13/viper"
	"github.com/tufanbarisyildirim/gonginx"
	"github.com/tufanbarisyildirim/gonginx/parser"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HostController struct {
}

func (x HostController) Init(ctx *gin.Context) {
	var hosts []model.VirtualHost
	var dc = x.updateVirtualHost(hosts)
	out, err := yaml.Marshal(dc)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": err.Error(), "data": nil})
		return
	}
	ctx.String(http.StatusOK, string(out))
	return
}

func (x HostController) List(ctx *gin.Context) {
	err := x.initVirtualHost()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  nil,
		"data": viper.Get("hosts").([]model.VirtualHost),
	})
}
func (x HostController) Show(ctx *gin.Context) {
	var id = ctx.Param("id")
	err := x.initVirtualHost()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	for _, h := range viper.Get("hosts").([]model.VirtualHost) {
		if h.Id == id {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  nil,
				"data": h,
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "未找到配置", "data": nil})
}
func (x HostController) Update(ctx *gin.Context) {
	var id = ctx.Param("id")
	var domain = ctx.PostForm("domain")
	var name = ctx.PostForm("name")
	var root = ctx.PostForm("root")
	var webRoot = ctx.PostForm("web_root")

	// 校验domain/root/web_root等参数格式
	if !strings.Contains(domain, ".") {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "域名格式错误", "data": domain})
		return
	}
	if strings.TrimSpace(root) == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "项目跟路径错误", "data": nil})
		return
	}

	err := x.initVirtualHost()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var findUpdate = false
	var hosts = viper.Get("hosts").([]model.VirtualHost)
	for i, h := range hosts {
		if h.Id == id {
			findUpdate = true
			h.Name = name
			h.Domain = domain
			h.Root = root
			h.WebRoot = webRoot
			hosts[i] = h
		}
	}
	if !findUpdate {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "修改配置存在", "data": nil})
		return
	}
	viper.Set("hosts", hosts)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "配置已经修改，请重启服务", "data": nil})
}
func (x HostController) Create(ctx *gin.Context) {
	var name = ctx.PostForm("name")
	var domain = ctx.PostForm("domain")
	var root = ctx.PostForm("root")
	var webRoot = ctx.PostForm("web_root")

	// 校验domain/root/web_root等参数格式
	if !strings.Contains(domain, ".") {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "域名格式错误", "data": domain})
		return
	}
	if strings.TrimSpace(root) == "" {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "项目跟路径错误", "data": nil})
		return
	}

	err := x.initVirtualHost()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var hosts = viper.Get("hosts").([]model.VirtualHost)
	for _, h := range hosts {
		if h.Domain == domain {
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "域名已经存在", "data": nil})
			return
		}
	}
	hosts = append(hosts, model.VirtualHost{
		Id:      x.hash(domain),
		Name:    name,
		Domain:  domain,
		Root:    root,
		WebRoot: webRoot,
		Port:    0,
	})
	viper.Set("hosts", hosts)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "配置已经添加，请重启服务", "data": nil})
}
func (x HostController) Delete(ctx *gin.Context) {
	var id = ctx.Param("id")
	err := x.initVirtualHost()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var newHosts []model.VirtualHost
	for _, h := range viper.Get("hosts").([]model.VirtualHost) {
		if h.Id == id {
			continue
		}
		newHosts = append(newHosts, h)
	}
	log.Println(x.toJson(newHosts))
	viper.Set("hosts", newHosts)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "已经删除配置", "data": nil})
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

func (x HostController) initVirtualHost() error {
	// 需要在配置docker前设置默认虚拟主机
	file, err := os.OpenFile("config.toml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		return err
	}
	err = viper.ReadConfig(file)
	if err != nil {
		log.Fatalln("[viper.ReadConfig error]", err.Error())
		return err
	}
	var hosts = []model.VirtualHost{
		{
			Id:      x.hash("default.me"),
			Name:    "default",
			Domain:  "default.me",
			Root:    filepath.Join(x.currentDirectory(), "dockerfile/nginx/html"),
			WebRoot: "",
			Port:    0,
		},
	}
	viper.SetConfigFile("config.toml")
	viper.SetDefault("READEME", "该配置自动生成，请勿修改")
	viper.SetDefault("app", "docker-lnmp")
	viper.SetDefault("hosts", hosts)
	err = viper.WriteConfig()
	if err != nil {
		log.Fatalln("[viper.WriteConfig error]", err.Error())
		return err
	}
	return nil
}

func (x HostController) hash(s string) string {
	var md5Hash = md5.New()
	_, err := io.WriteString(md5Hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5Hash.Sum(nil))[:8]
}
