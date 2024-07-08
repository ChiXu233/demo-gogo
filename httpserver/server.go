package httpserver

import (
	"demo-gogo/api/handler"
	"demo-gogo/config"
	"demo-gogo/httpserver/app"
	"demo-gogo/httpserver/middleware"
	"github.com/gin-gonic/gin"
)

const (
	ApiCallBack = "callback"
	ApiMap      = "map"
	ApiRobot    = ""
	ApiDebug    = "debug"
	ApiVersion  = ""
)

func CreateHttpServer() *gin.Engine {
	gin.SetMode(config.Conf.APP.Mode)
	engine := gin.New()
	middlewareList := []gin.HandlerFunc{
		gin.Logger(),
		// 日志组件增强，用来打印gin的入参
		middleware.RequestInfo(),
		gin.Recovery(),
	}
	// 路由注册，中间件引入
	RegisterRoutes(engine, middlewareList)
	return engine
}

func RegisterRoutes(router *gin.Engine, middlewares []gin.HandlerFunc) {
	// 为全局路由注册中间件
	router.Use(middlewares...)
	// 捕捉不允许的方法
	router.NoMethod(app.MethodNotFound)
	router.NoRoute(app.HandleNotFound)
	// 静态路由
	router.Static("/files", "./files")

	// 设置系统路径上下文
	contextPath := router.Group(config.Conf.APP.ContextPath)

	v1 := contextPath.Group(ApiVersion)
	//// api接口注册鉴权中间件
	//v1.Use(middleware.Auth())
	restHandler := handler.NewHandler()
	v1.Group("")
	{
		v1.GET("/ping", restHandler.V1Ping)
	}
	callBack := contextPath.Group(ApiCallBack)
	callBack.Group("")

	robot := contextPath.Group(ApiRobot)
	robot.Group("")

	m := contextPath.Group(ApiMap)
	//m.Use(middleware.Auth())
	{
		m.POST("/mapping", restHandler.CreateOrUpdateMap)
		m.GET("/maps", restHandler.ListMap)
		m.DELETE("/mapping/:id", restHandler.DeleteMap)

		m.POST("/map_info", restHandler.CreateOrUpdateMapInfo)
		m.GET("/map_info", restHandler.ListMapInfosInfo)
		m.DELETE("/map_info/:id", restHandler.DeleteMapInfo)
		m.POST("/map_info_nodes:info_id", restHandler.CreateOrUpdateNode) //生成路径节点
		m.GET("/map_info_nodes", restHandler.ListMapNodes)
		m.DELETE("/map_info_nodes/:id", restHandler.DeleteMapNodes)
		m.POST("/map_info_routes/:info_id", restHandler.CreateOrUpdateMapRoutes) //生成路径节点+路径
		m.GET("/map_info_routes", restHandler.ListMapRoutes)                     //查找路径
		m.DELETE("/map_info_routes/:id", restHandler.DeleteMapRoute)
		m.POST("/check_route", restHandler.CheckRoute) //检验路径
		m.GET("/map_infos/:map_id", restHandler.ListMapInfo)
		m.POST("/map_nodes_batch/", restHandler.BatchDeleteMapNodes) //批量删除路径节点

	}

	if config.Conf.APP.Mode == gin.DebugMode {
		debug := contextPath.Group(ApiDebug)
		debug.Group("")
	}
}
