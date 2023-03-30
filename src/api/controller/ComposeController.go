package controller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/util"
	"net/http"
)

type ComposeController struct {
}

// docker-compose.yaml 文件路径
// 状态查看
// 		当前compose中

// docker container ls 列出容器，遍历出 Label.com.docker.compose.project 同名的容器列表 即为当前分组

func (x ComposeController) connect(ctx *gin.Context) *client.Client {
	_clientInstance, err := util.ConnectDocker()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return nil
	}
	return _clientInstance
}

func (x ComposeController) Status(ctx *gin.Context) {
	dClient := x.connect(ctx)
	listSummary, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var resultList []any
	for _, item := range listSummary {
		if _, ok := item.Labels["com.docker.compose.project"]; !ok {
			continue
		}
		if item.Labels["com.docker.compose.project"] != config.AppName {
			continue
		}

		var ports []string
		for _, p := range item.Ports {
			if p.PublicPort == 0 {
				ports = append(ports, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
			} else {
				ports = append(ports, fmt.Sprintf("%d:%d", p.PublicPort, p.PrivatePort))
			}
		}

		resultList = append(resultList, gin.H{
			"id":       item.ID[:8],
			"name":     item.Names[0],
			"image":    item.Image,
			"image_id": item.ImageID,
			"labels": gin.H{
				"project":     item.Labels["com.docker.compose.project"],
				"service":     item.Labels["com.docker.compose.service"],
				"version":     item.Labels["com.docker.compose.version"],
				"config_file": item.Labels["com.docker.compose.project.config_files"],
				"working_dir": item.Labels["com.docker.compose.project.working_dir"],
			},
			"state":      item.State,
			"ports":      ports,
			"status":     item.Status,
			"created_at": item.Created,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": resultList})
	return
}

func (x ComposeController) Images(ctx *gin.Context) {
	dClient := x.connect(ctx)
	imagesListSummary, err := dClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": imagesListSummary})
	return
}

func (x ComposeController) Containers(ctx *gin.Context) {
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

func (x ComposeController) Stop(ctx *gin.Context) {
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
