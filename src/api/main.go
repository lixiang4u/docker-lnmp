package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/controller"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run web server",
				Action: func(ctx *cli.Context) error {
					runServer()
					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init environment",
				Action: func(ctx *cli.Context) error {
					//
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
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

	r.GET("/host/init", hostController.Init)
	r.GET("/host/list", hostController.List)
	r.GET("/host/show/:id", hostController.Show)
	r.POST("/host/create", hostController.Create)
	r.PUT("/host/update/:id", hostController.Update)
	r.DELETE("/host/delete/:id", hostController.Delete)

	r.GET("/image/list", imageController.List)
	r.DELETE("/image/remove/:id", imageController.Remove)

	r.GET("/container/list", containerController.List)
	r.POST("/container/start/:id", containerController.Start)
	r.POST("/container/stop/:id", containerController.Stop)
	r.POST("/container/remove/:id", containerController.Remove)
	r.GET("/container/status/:id", containerController.Status)
	r.GET("/container/logs/:id", containerController.Logs)

	r.GET("/project/list", projectController.List)
	r.POST("/project/start", projectController.Start)
	r.POST("/project/stop", projectController.Stop)
	r.POST("/project/remove", projectController.Remove)
	r.POST("/project/rebuild", projectController.Rebuild)

	_ = r.Run(":8086")
}
