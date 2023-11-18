package main

import (
	"github.com/chenchenyu/gomysqldemo/api"
	"github.com/chenchenyu/gomysqldemo/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	dao.Init()
	g := gin.Default()
	api.InitRouter(g)
	g.Run(":9000")
}
