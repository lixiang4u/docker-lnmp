package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/config"
	"github.com/lixiang4u/docker-lnmp/controller"
	"github.com/lixiang4u/docker-lnmp/model"
	"github.com/lixiang4u/docker-lnmp/util"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run web server",
				Action: func(ctx *cli.Context) error {
					if err := checkService(); err != nil {
						return err
					}
					runServer()
					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init environment",
				Action: func(ctx *cli.Context) error {
					if err := checkService(); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func checkService() error {
	log.Println("check docker-compose.yaml...")

	// 连接docker服务
	dClient, err := util.ConnectDocker()
	if err != nil {
		return err
	}
	defer func() { _ = dClient.Close() }()

	// 计算服务是否存在
	listSummary, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	var resultList []model.Container
	for _, item := range listSummary {
		if item.Labels["com.docker.compose.project"] == "" {
			continue
		}
		if item.Labels["com.docker.compose.project"] != config.AppName {
			continue
		}
		resultList = append(resultList, model.Container{Id: item.ID, Name: item.Names[0]})
	}
	configFile := filepath.Join(util.AppDirectory(), "docker-compose.yaml")
	log.Println("[configFile]", configFile)
	if len(resultList) == 0 {
		// 生成 docker-compose.yaml
		err := new(controller.ComposeController).GenerateConfig(configFile)
		if err != nil {
			return err
		}

		// 运行docker compose up
		log.Println(fmt.Sprintf("正在构建%s，稍等片刻...", config.AppName))
		_, err = model.ComposeDownUp(configFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func runServer() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	var hostController = new(controller.HostController)
	var containerController = new(controller.ContainerController)
	var imageController = new(controller.ImageController)
	var projectController = new(controller.ProjectController)

	r.NoRoute(func(ctx *gin.Context) { ctx.Redirect(http.StatusPermanentRedirect, "/") })
	r.StaticFile("/", filepath.Join(util.AppDirectory(), "src/web/dist/index.html"))
	r.Static("/assets", filepath.Join(util.AppDirectory(), "src/web/dist/assets"))

	api := r.Group("/api")
	api.GET("/host/init", hostController.Init)
	api.GET("/host/list", hostController.List)
	api.GET("/host/show/:id", hostController.Show)
	api.POST("/host/create", hostController.Create)
	api.PUT("/host/update/:id", hostController.Update)
	api.DELETE("/host/delete/:id", hostController.Delete)

	api.GET("/image/list", imageController.List)
	api.DELETE("/image/remove/:id", imageController.Remove)
	api.GET("/image/run/:id", imageController.Run)

	api.GET("/container/list", containerController.List)
	api.POST("/container/start/:id", containerController.Start)
	api.POST("/container/stop/:id", containerController.Stop)
	api.POST("/container/remove/:id", containerController.Remove)
	api.GET("/container/status/:id", containerController.Status)
	api.GET("/container/logs/:id", containerController.Logs)

	api.GET("/project/list", projectController.List)
	api.POST("/project/start", projectController.Start)
	api.POST("/project/stop", projectController.Stop)
	api.POST("/project/remove", projectController.Remove)
	api.POST("/project/rebuild", projectController.Rebuild)

	_ = r.Run(":8086")
}
