package main

import (
	"ginchat/config"
	"ginchat/router"
)

func main() {
	//加载数据库
	config.Init()
	router := router.NewRouter()
	_ = router.Run(":8080")

}
