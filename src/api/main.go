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
	var clientController = new(controller.DockerClientController)
	var composeController = new(controller.ComposeController)
	r.GET("/host/init", hostController.Init)
	r.GET("/host/list", hostController.List)
	r.GET("/host/show/:id", hostController.Show)
	r.POST("/host/create", hostController.Create)
	r.PUT("/host/update/:id", hostController.Update)
	r.DELETE("/host/delete/:id", hostController.Delete)
	r.GET("/docker/images", clientController.Images)
	r.GET("/docker/containers", clientController.Containers)
	r.GET("/docker/container/stop", clientController.Stop)
	r.GET("/compose/status", composeController.Status)

	_ = r.Run(":8086")
}
