package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/controller"
)

func main() {
	r := gin.Default()

	var hostController = new(controller.HostController)
	r.GET("/host/init", hostController.Init)
	r.GET("/host/list", hostController.List)
	r.GET("/host/show/:id", hostController.Show)
	r.POST("/host/create", hostController.Create)
	r.PUT("/host/update/:id", hostController.Update)
	r.DELETE("/host/delete/:id", hostController.Delete)

	_ = r.Run(":8086")
}
