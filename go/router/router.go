package router

import (
	controller2 "SkyLine/controller"
	_ "SkyLine/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() {
	r := gin.Default()

	gin.SetMode(viper.GetString("server.mode"))
	InitBasePlatformRouter(r)
	setPort := viper.GetString("server.port")

	err := r.Run(fmt.Sprintf(":%s", setPort))
	if err != nil {
		fmt.Println("路由初始化失败")
	}
}

// InitBasePlatformRouter 该方法中添加路由
func InitBasePlatformRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	// ===========================swagger文档路由===================================
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//=============================================================================

	apiRouter := r.Group("/douyin")

	// ============================swagger文档创建测试===============================
	//get请求
	apiRouter.GET("/swaggereget/", controller2.Swaggerget)
	//post请求
	apiRouter.GET("/swaggerpost/", controller2.Swaggerpost)
	// ============================================================================

	// basic apis
	apiRouter.GET("/feed/", controller2.Feed)
	apiRouter.GET("/user/", controller2.UserInfo)
	apiRouter.POST("/user/register/", controller2.Register)
	apiRouter.POST("/user/login/", controller2.Login)
	apiRouter.POST("/publish/action/", controller2.Publish)
	apiRouter.GET("/publish/list/", controller2.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller2.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller2.FavoriteList)
	apiRouter.POST("/comment/action/", controller2.CommentAction)
	apiRouter.GET("/comment/list/", controller2.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller2.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller2.FollowList)
	apiRouter.GET("/relation/follower/list/", controller2.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller2.FriendList)
	apiRouter.GET("/message/chat/", controller2.MessageChat)
	apiRouter.POST("/message/action/", controller2.MessageAction)

}
