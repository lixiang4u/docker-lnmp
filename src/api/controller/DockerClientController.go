package controller

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DockerClientController struct {
}

var _clientInstance *client.Client

func (x DockerClientController) connect(ctx *gin.Context) *client.Client {
	if _clientInstance != nil {
		return _clientInstance
	}

	var err error
	// Error response from daemon: client version 1.42 is too new. Maximum supported API version is 1.41
	_clientInstance, err = client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return nil
	}
	return _clientInstance
}

func (x DockerClientController) Containers(ctx *gin.Context) {
	dClient := x.connect(ctx)
	containerListSummary, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": containerListSummary})
	return
}

func (x DockerClientController) Stop(ctx *gin.Context) {
	var id = ctx.Query("id")
	dClient := x.connect(ctx)
	err := dClient.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "容器已经停止", "data": nil})
	return
}
