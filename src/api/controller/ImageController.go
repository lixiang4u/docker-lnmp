package controller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/util"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type ImageController struct {
}

func (x ImageController) connect(ctx *gin.Context) *client.Client {
	_clientInstance, err := util.ConnectDocker()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return nil
	}
	return _clientInstance
}

func (x ImageController) List(ctx *gin.Context) {
	var project = ctx.Query("project")
	dClient := x.connect(ctx)
	listSummary, err := dClient.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	var resultList []any
	for _, item := range listSummary {
		if project != "all" {
			if _, ok := item.Labels["com.docker.compose.project"]; !ok {
				continue
			}
			if item.Labels["com.docker.compose.project"] != config.AppName {
				continue
			}
		}
		resultList = append(resultList, gin.H{
			"id":  x.shortId(item.ID),
			"tag": item.RepoTags[0],
			"labels": gin.H{
				"project": item.Labels["com.docker.compose.project"],
				"service": item.Labels["com.docker.compose.service"],
				"version": item.Labels["com.docker.compose.version"],
			},
			"size":         item.Size,
			"virtual_size": item.VirtualSize,
			"created_at":   item.Created,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": resultList})
	return
}

func (x ImageController) Remove(ctx *gin.Context) {
	var id = ctx.Param("id")
	dClient := x.connect(ctx)
	_, err := dClient.ImageRemove(context.Background(), id, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "镜像已经移除", "data": nil})
	return
}

func (x ImageController) Run(ctx *gin.Context) {
	var id = ctx.Param("id")

	var idLen = len(id)
	if idLen < 2 {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "镜像ID错误", "data": nil})
		return
	}

	dClient := x.connect(ctx)
	listSummary, err := dClient.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	for _, tmpImage := range listSummary {
		log.Println("======> ", x.shortId(tmpImage.ID), id)
		if x.shortId(tmpImage.ID) == id {
			tmpName := fmt.Sprintf(
				"%s-%d",
				strings.ReplaceAll(tmpImage.RepoTags[0], ":", "-"),
				1000+rand.Intn(999),
			)
			response, err := dClient.ContainerCreate(
				context.Background(),
				&container.Config{
					Image: tmpImage.RepoTags[0],
				},
				&container.HostConfig{
					// PublishAllPorts
					//PortBindings: map[nat.Port][]nat.PortBinding{
					//	"80": []nat.PortBinding{
					//		{
					//			HostIP:   "0.0.0.0",
					//			HostPort: "80",
					//		},
					//	},
					//},
				},
				&network.NetworkingConfig{},
				&v1.Platform{},
				tmpName,
			)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error(), "data": nil})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": response, "tmpName": tmpName})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": "镜像不存在", "data": nil})
	return
	//		_, err := c.ContainerCreate(ctx, &containertypes.Config{Image: "busybox:latest"}, &containertypes.HostConfig{}, nil, &p, "")
}

func (x ImageController) shortId(id string) string {
	return id[7 : 7+8]
}
