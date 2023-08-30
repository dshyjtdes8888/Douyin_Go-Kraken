package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"simple-demo/data"
	"strconv"
	"time"
  "os"
)

// FavoriteAction 函数响应点赞，1为点赞，2为取消点赞。
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	vid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actiontype, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)

	// 检查用户是否存在，如果用户不存在，则返回StatusCode为1，表示用户不存在。
	var user User
	result := data.Db.Where("Name = ?", token).First(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	// 开始事务
	tx := data.Db.Begin()

	if actiontype == 1 {
		// 更新用户的favorite_count和视频的喜欢状态和计数
		if err := tx.Model(&User{Id: user.Id}).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 2, StatusMsg: "Failed to update user favorite count"})
			return
		}
		if err := tx.Model(&Video{Id: vid}).Update("is_favorite", true).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 3, StatusMsg: "Failed to update video favorite status"})
			return
		}
		if err := tx.Model(&Video{Id: vid}).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 4, StatusMsg: "Failed to update video favorite count"})
			return
		}

		// 添加到favorite表
		favorite := Favorite{
			UserId:  user.Id,
			VideoId: int(vid),
		}
		if err := tx.Create(&favorite).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 5, StatusMsg: "Failed to create favorite entry"})
			return
		}

		// 提交事务
		tx.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		// 处理取消点赞操作
		if user.FavoriteCount > 0 {
			if err := tx.Model(&User{Id: user.Id}).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, Response{StatusCode: 6, StatusMsg: "Failed to update user favorite count"})
				return
			}
			if err := tx.Model(&Video{Id: vid}).Update("is_favorite", false).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, Response{StatusCode: 7, StatusMsg: "Failed to update video favorite status"})
				return
			}

			// 从favorite表中删除对应的条目
			if err := tx.Where("user_id = ? AND video_id = ?", user.Id, vid).Delete(Favorite{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, Response{StatusCode: 8, StatusMsg: "Failed to delete favorite entry"})
				return
			}

			// 更新视频的favorite_count
			var video Video
			if err := tx.Where("id = ?", vid).Find(&video).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, Response{StatusCode: 9, StatusMsg: "Failed to find video"})
				return
			}
			if video.FavoriteCount > 0 {
				if err := tx.Model(&Video{Id: vid}).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, Response{StatusCode: 10, StatusMsg: "Failed to update video favorite count"})
					return
				}
			}

			// 提交事务
			tx.Commit()
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		}
	}
}

// FavoriteList 函数返回所有用户的点赞视频列表。
func FavoriteList(c *gin.Context) {
	// 更改所有视频的IsFavorite字段为false
	if err := data.Db.Model(&Video{}).Where("1 = 1").Update("is_favorite", false).Error; err != nil {
    c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Failed to update videos"})
		return
	}

	// 初始化一个切片来存储用户已收藏的视频ID
	var favoriteVideoIDs []int

	// 获取用户的身份信息，例如令牌
	token := c.Query("token")
	userid := c.Query("user_id")

	if token != "" {
		// 根据令牌查找用户
		var user User
		result := data.Db.Where("id = ?", userid).First(&user)
		if result.Error != nil {
    c.JSON(http.StatusInternalServerError, Response{StatusCode: 2, StatusMsg: "User doesn't exist"})
			return
		}

		// 查询用户的收藏视频ID
		var favorites []*Favorite
		result = data.Db.Where("user_id = ?", user.Id).Find(&favorites)
		if result.Error != nil {
          c.JSON(http.StatusInternalServerError, Response{StatusCode: 3, StatusMsg: "Failed to find favorites"})
			return
		}

		// 收集用户已收藏的视频ID
		for _, favorite := range favorites {
			favoriteVideoIDs = append(favoriteVideoIDs, favorite.VideoId)
		}

		// 更新用户已收藏的视频的 is_favorite 字段
		if len(favoriteVideoIDs) > 0 {
			if err := data.Db.Model(&Video{}).Where("id IN (?)", favoriteVideoIDs).Update("is_favorite", true).Error; err != nil {
            c.JSON(http.StatusInternalServerError, Response{StatusCode: 4, StatusMsg: "Failed to update videos"})
				return
			}
		}

		var videos []*Video
		result = data.Db.Where("is_favorite = ?", true).Find(&videos)
		if result.Error != nil {
          c.JSON(http.StatusInternalServerError, Response{StatusCode: 5, StatusMsg: "Failed to fetch videos"})
			return
		}

      // 从环境变量获取动态域名信息
	paasURL := os.Getenv("paas_url")

	if paasURL == "" {
    c.JSON(http.StatusInternalServerError, Response{StatusCode: 6, StatusMsg: "Failed to get pass_url"})
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

		c.JSON(http.StatusOK, FeedResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			VideoList:  videoList,
			NextTime:   time.Now().Unix(),
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 6, StatusMsg: "User doesn't login"})
	}
}
