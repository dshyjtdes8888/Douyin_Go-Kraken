package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-demo/data"
	"time"
  "os"
  "fmt"
)

func Feed(c *gin.Context) {
	var videos []*Video
	result := data.Db.Preload("Author").Find(&videos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "Failed to fetch videos"})
		return
	}

      	// 获取动态域名信息从环境变量
	paasURL := os.Getenv("paas_url")

	if paasURL == "" {
		fmt.Println("环境变量 paas_url 未设置")
		return
	}

  

	var videoList []VideoListItem // 自定义的结构体，用于映射视频列表项

	for _, video := range videos {
    video.Author.Avatar = "https://"+paasURL+"/static/"+video.Author.Avatar
    video.Author.BackgroundImage = "https://"+paasURL+"/static/"+video.Author.BackgroundImage
    
		videoListItem := VideoListItem{
			ID:            video.Id,
			Author:        video.Author,
			PlayURL:       "https://"+paasURL+"/static/"+video.PlayUrl,
			CoverURL:      "https://"+paasURL+"/static/"+video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		videoList = append(videoList, videoListItem)
	}

	feedResponse := FeedResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		NextTime:   time.Now().Unix(),
		VideoList:  videoList,
	}

	c.JSON(http.StatusOK, feedResponse)
}
