package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/util"
	"io"
	"log"
	"net/http"
)

type ContainerController struct {
}

func (x ContainerController) connect(ctx *gin.Context) *client.Client {
	_clientInstance, err := util.ConnectDocker()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return nil
	}
	return _clientInstance
}

func (x ContainerController) List(ctx *gin.Context) {
	var imageId = ctx.Query("image_id")

	newArgs := filters.NewArgs()
	if imageId != "" {
		newArgs.Add("ancestor", imageId)
	}

	dClient := x.connect(ctx)
	listSummary, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: newArgs,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var resultList []any
	for _, item := range listSummary {
		if _, ok := item.Labels["com.docker.compose.project"]; !ok {
			//continue
		}
		if item.Labels["com.docker.compose.project"] != config.AppName {
			//continue
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

func (x ContainerController) Start(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	err := dClient.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "容器已经启动", "data": nil})
	return
}

func (x ContainerController) Stop(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	err := dClient.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "容器已经停止", "data": nil})
	return
}

func (x ContainerController) Status(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	stats, err := dClient.ContainerStatsOneShot(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	defer func() { _ = stats.Body.Close() }()

	log.Printf("[stats] %v", stats)

	v := types.Stats{}
	err = json.NewDecoder(stats.Body).Decode(&v)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": stats, "v": v})
	return
}

func (x ContainerController) Remove(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	err := dClient.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "容器已经移除", "data": nil})
	return
}

func (x ContainerController) Logs(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	reader, err := dClient.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{
		Tail:       "2000",
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
		Timestamps: true,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	defer func() { _ = reader.Close() }()
	bs, err := io.ReadAll(reader)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "日志查询",
		"data": util.DockerLogFormat(string(bs)),
	})
	return
}
