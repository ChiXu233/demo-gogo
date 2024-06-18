package main

import (
	"demo-gogo/config"
	"demo-gogo/database"
	"demo-gogo/httpserver"
	"demo-gogo/utils/redis"
	"fmt"
	log "github.com/wonderivan/logger"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic("load config with error:" + err.Error())
	}

	err = config.InitLog()
	if err != nil {
		panic("init log with error:" + err.Error())
	}

	err = database.InitDB()
	if err != nil {
		panic("init database with error:" + err.Error())
	}

	//if config.Conf.OSS.Type == config.CONF_OSS_MINIO {
	//	err = storage.InitStorage()
	//	if err != nil {
	//		panic("init storage with error:" + err.Error())
	//	}
	//}
	//err = env.InitPython()
	//if err != nil {
	//	panic("init pythonEnv with error:" + err.Error())
	//}

	err = redis.InitRedis()
	if err != nil {
		panic("init redis with error:" + err.Error())
	}

	//if config.Conf.Compute.PullInterval > 0 {
	//	go PullCalculateProgress()
	//}

	server := httpserver.CreateHttpServer()
	listenAddress := fmt.Sprintf("0.0.0.0:%d", config.Conf.APP.Port)

	if err = server.Run(listenAddress); err != nil {
		log.Error("ma_teach exit with error: %v", err)
	}

}
