package controller

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/util"
	"net/http"
)

type ProjectController struct {
	ContainerController
}

func (x ProjectController) connect(ctx *gin.Context) *client.Client {
	_clientInstance, err := util.ConnectDocker()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return nil
	}
	return _clientInstance
}

func (x ProjectController) List(ctx *gin.Context) {
	var imageId = ctx.Query("image_id")
	var project = ctx.Query("project")

	listSummary, err := x.findContainers(x.connect(ctx), imageId, project)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": x.containersToProjects(listSummary),
	})
	return
}

func (x ProjectController) Start(ctx *gin.Context) {
	var projectName = ctx.Param("projectName")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": projectName})
	return
}

func (x ProjectController) Stop(ctx *gin.Context) {
	var projectName = ctx.Param("projectName")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": projectName})
	return
}

func (x ProjectController) Remove(ctx *gin.Context) {
	var projectName = ctx.Param("projectName")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": projectName})
	return
}
