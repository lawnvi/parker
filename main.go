package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"parker/app/part_test/controller"
	"parker/config/database"
	"parker/config/log"
	"parker/config/routers"
)

func main() {
	//加载log
	//log.Init(log.NoWrite, "./parker.log")

	gin.SetMode(gin.DebugMode)
	//加载数据库配置
	if err := database.Init(); err != nil {
		log.E("init", fmt.Sprintf("database init error, %v", err))
		return
	}
	defer func() { _ = database.DB.Close() }()

	//加载多个app的路由配置
	routers.Include(controller.UserRouters)
	//初始化路由
	r := routers.Init()
	if err := r.Run(":8088"); err != nil {
		log.E("init", fmt.Sprintf("start service failed, %v", err))
	}
	log.I("init", "web start success")
}
