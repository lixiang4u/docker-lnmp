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
	"github.com/lixiang4u/docker-lnmp/model"
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

	listSummary, err := x.findContainers(x.connect(ctx), imageId, "")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": listSummary})
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

func (x ContainerController) findContainers(dockerClient *client.Client, imageId, projectName string) ([]model.Container, error) {
	newArgs := filters.NewArgs()
	if imageId != "" {
		newArgs.Add("ancestor", imageId)
	}

	listSummary, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: newArgs})
	if err != nil {
		return nil, err
	}
	var resultList []model.Container
	for _, item := range listSummary {
		if projectName != "" && item.Labels["com.docker.compose.project"] != projectName {
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
		resultList = append(resultList, model.Container{
			Id:      item.ID,
			Name:    item.Names[0],
			Image:   item.Image,
			ImageId: item.ImageID,
			Labels: model.Label{
				Project:    item.Labels["com.docker.compose.project"],
				Service:    item.Labels["com.docker.compose.service"],
				Version:    item.Labels["com.docker.compose.version"],
				ConfigFile: item.Labels["com.docker.compose.project.config_files"],
				WorkingDir: item.Labels["com.docker.compose.project.working_dir"],
			},
			State:     item.State,
			Ports:     ports,
			Status:    item.Status,
			CreatedAt: item.Created,
		})
	}
	return resultList, err
}

func (x ContainerController) containersToProjects(containerList []model.Container) map[string][]model.Container {
	var resultMap = make(map[string][]model.Container)
	for _, item := range containerList {
		if item.Labels.Project == "" {
			resultMap[item.Name] = append(resultMap[item.Name], item)
		} else {
			resultMap[item.Labels.Project] = append(resultMap[item.Labels.Project], item)
		}
	}
	return resultMap
}
