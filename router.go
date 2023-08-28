package main

import (
  "simple-demo/controller"
  "github.com/gin-gonic/gin"
  "net/http"
)

func initRouter(r *gin.Engine) {
  //public directory is used to serve static resources
  r.Static("/static", "./public")
  r.LoadHTMLGlob("templates/*")

  //home page
  r.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
  })
  
  apiRouter := r.Group("/douyin")

  // basic apis路由定义
	apiRouter.GET("/feed/", controller.Feed)                // 获取动态信息的API
	apiRouter.GET("/user/", controller.UserInfo)            // 获取用户信息的API
	apiRouter.POST("/user/register/", controller.Register)  // 用户注册的API
	apiRouter.POST("/user/login/", controller.Login)        // 用户登录的API
	apiRouter.POST("/publish/action/", controller.Publish)  // 发布动态的API
	apiRouter.GET("/publish/list/", controller.PublishList) // 获取动态列表的API

	//额外APIs - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction) // 点赞/取消点赞动态的API
	apiRouter.GET("/favorite/list/", controller.FavoriteList)      // 获取用户点赞过的动态列表的API
	apiRouter.POST("/comment/action/", controller.CommentAction)   // 发布评论的API
	apiRouter.GET("/comment/list/", controller.CommentList)        // 获取动态的评论列表的API

	//额外APIs - II
	apiRouter.POST("/relation/action/", controller.RelationAction)     // 关注/取消关注用户的API
	apiRouter.GET("/relation/follow/list/", controller.FollowList)     // 获取用户关注列表的API
	apiRouter.GET("/relation/follower/list/", controller.FollowerList) // 获取用户粉丝列表的API
	apiRouter.GET("/relation/friend/list/", controller.FriendList)     // 获取用户好友列表的API
	apiRouter.GET("/message/chat/", controller.MessageChat)            // 获取用户的私信聊天记录的API
	apiRouter.POST("/message/action/", controller.MessageAction)       // 发送私信的API
}
