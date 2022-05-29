package main

import (
	"log"

	"github.com/RxDAF/Env/cfg"
	"github.com/RxDAF/Env/server"
)

func main() {
	cfg, err := cfg.NewCFG("config.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	// 与主服务器建立连接，拉取二进制文件
	server := server.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Println(err)
		return
	}
}
