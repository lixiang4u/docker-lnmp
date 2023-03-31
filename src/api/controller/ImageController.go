package controller

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/util"
	"net/http"
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
		if _, ok := item.Labels["com.docker.compose.project"]; !ok {
			//continue
		}
		if item.Labels["com.docker.compose.project"] != config.AppName {
			//continue
		}
		resultList = append(resultList, gin.H{
			"id":  item.ID[7 : 7+8],
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
