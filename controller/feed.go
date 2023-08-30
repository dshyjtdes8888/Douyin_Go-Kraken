package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-demo/data"
	"time"
  "os"
)

// FeedResponse 自定义结构体用于接口响应
type FeedResponse struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	NextTime   int64           `json:"next_time"`
	VideoList  []VideoListItem `json:"video_list"`
}

func Feed(c *gin.Context) {
	var videos []*Video
	result := data.Db.Preload("Author").Find(&videos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "Failed to fetch videos"})
		return
	}

	token := c.Query("token")

	// 更改所有视频的IsFavorite字段为false
	if err := data.Db.Model(&Video{}).Where("1 = 1").Update("is_favorite", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "Failed to update videos"})
		return
	}

	// 初始化一个切片来存储用户已点赞的视频Id
	var favoriteVideoIDs []int

	// 检查用户是否已登录并且具有有效的token参数
	if token != "" {
		var user User
		result = data.Db.Where("name = ?", token).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "无法找到用户"})
			return
		}

		// 查询用户的点赞视频
		var favorites []*Favorite
		result = data.Db.Where("user_id = ?", user.Id).Find(&favorites)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "无法查询收藏记录"})
			return
		}

		// 获取用户已点赞的视频ID
		for _, favorite := range favorites {
			favoriteVideoIDs = append(favoriteVideoIDs, favorite.VideoId)
		}

		// 更新用户已收藏的视频的 is_favorite 字段
		if len(favoriteVideoIDs) > 0 {
			if err := data.Db.Model(&Video{}).Where("id IN (?)", favoriteVideoIDs).Update("is_favorite", true).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "status_msg": "无法更新收藏的视频"})
				return
			}
		}
	}

  // 从环境变量获取动态域名信息
	paasURL := os.Getenv("paas_url")

	if paasURL == "" {
		    c.JSON(http.StatusOK, Response{StatusCode: 6, StatusMsg: "Failed to get pass_url"})
		return
	}

	var videoList []VideoListItem

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
		StatusMsg:  "success",
		NextTime:   time.Now().Unix(),
		VideoList:  videoList,
	}

	c.JSON(http.StatusOK, feedResponse)
}
