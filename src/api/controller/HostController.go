package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/model"
	"github.com/lixiang4u/docker-lnmp/util"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
)

type HostController struct {
}

func (x HostController) Init(ctx *gin.Context) {
	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var hosts = model.ConvertToVirtualHost(viper.Get("hosts"))

	var dc = model.UpdateVirtualHost(hosts)
	out, err := yaml.Marshal(dc)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": err.Error(), "data": nil})
		return
	}
	ctx.String(http.StatusOK, string(out))
	return
}

func (x HostController) List(ctx *gin.Context) {
	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  nil,
		"data": model.ConvertToVirtualHost(viper.Get("hosts")),
	})
}

func (x HostController) Show(ctx *gin.Context) {
	var id = ctx.Param("id")
	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	for _, h := range model.ConvertToVirtualHost(viper.Get("hosts")) {
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

	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var findUpdate = false
	var hosts = model.ConvertToVirtualHost(viper.Get("hosts"))
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
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "配置已经修改，请【重新构建】容器", "data": nil})
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

	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var hosts = model.ConvertToVirtualHost(viper.Get("hosts"))
	for _, h := range hosts {
		if h.Domain == domain {
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "域名已经存在", "data": nil})
			return
		}
	}
	hosts = append(hosts, model.VirtualHost{
		Id:      util.StringHash(domain),
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
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "配置已经添加，请【重新构建】容器", "data": nil})
}

func (x HostController) Delete(ctx *gin.Context) {
	var id = ctx.Param("id")
	err := model.InitVirtualHostConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var newHosts []model.VirtualHost
	var oldHosts = model.ConvertToVirtualHost(viper.Get("hosts"))
	for _, h := range oldHosts {
		if h.Id == id {
			continue
		}
		newHosts = append(newHosts, h)
	}
	viper.Set("hosts", newHosts)
	err = viper.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "已经删除配置", "data": nil})
}
