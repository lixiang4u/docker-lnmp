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

	r.StaticFile("/", "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\src\\web\\dist\\index.html")
	r.Static("/assets", "D:\\repo\\github.com\\lixiang4u\\docker-lnmp\\src\\web\\dist\\assets")

	api := r.Group("/api")
	api.GET("/host/init", hostController.Init)
	api.GET("/host/list", hostController.List)
	api.GET("/host/show/:id", hostController.Show)
	api.POST("/host/create", hostController.Create)
	api.PUT("/host/update/:id", hostController.Update)
	api.DELETE("/host/delete/:id", hostController.Delete)

	api.GET("/image/list", imageController.List)
	api.DELETE("/image/remove/:id", imageController.Remove)

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
