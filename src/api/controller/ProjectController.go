package controller

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/model"
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
	var containerId = ctx.PostForm("container_id")
	var projectName = ctx.PostForm("project_name")

	var dClient = x.connect(ctx)
	listSummary, err := x.findContainersByProjectOrId(x.connect(ctx), projectName, containerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	for _, tmpContainer := range listSummary {
		err := dClient.ContainerStart(context.Background(), tmpContainer.Id, types.ContainerStartOptions{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": nil})
	return
}

func (x ProjectController) Stop(ctx *gin.Context) {
	var containerId = ctx.PostForm("container_id")
	var projectName = ctx.PostForm("project_name")

	var dClient = x.connect(ctx)
	listSummary, err := x.findContainersByProjectOrId(dClient, projectName, containerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	for _, tmpContainer := range listSummary {
		err := dClient.ContainerStop(context.Background(), tmpContainer.Id, container.StopOptions{})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": nil})
	return
}

func (x ProjectController) Remove(ctx *gin.Context) {
	var containerId = ctx.PostForm("container_id")
	var projectName = ctx.PostForm("project_name")

	var dClient = x.connect(ctx)
	listSummary, err := x.findContainersByProjectOrId(dClient, projectName, containerId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	for _, tmpContainer := range listSummary {
		err := dClient.ContainerRemove(context.Background(), tmpContainer.Id, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": nil})
	return
}

// 根据项目名或者容器ID查找容器列表（项目名优先）
func (x ContainerController) findContainersByProjectOrId(dockerClient *client.Client, projectName, containerId string) ([]model.Container, error) {
	if projectName != "" {
		return x.findContainers(dockerClient, "", projectName)
	}

	listSummary, err := x.findContainers(dockerClient, "", "")
	if err != nil {
		return nil, err
	}
	var resultList []model.Container
	for _, item := range listSummary {
		if item.Id != containerId {
			continue
		}
		resultList = append(resultList, item)
	}
	return resultList, nil
}
