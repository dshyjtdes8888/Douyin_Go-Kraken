package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"simple-demo/data"
  "os"
  "fmt"
)

// usersLoginInfo是一个map，用于存储用户信息，以用户名+密码为键。
// 该变量中的用户数据在每次服务器启动时会被清空，仅用于演示。
// 测试数据：用户名=zhanglei，密码=douyin

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//// userIdSequence是一个原子整数，用于生成用户的唯一ID。
//var userIdSequence = int64(1)

// UserLoginResponse 结构体表示用户登录的响应数据。
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// UserResponse 结构体表示用户信息的响应数据。
type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register函数用于处理用户注册的请求。
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var user User
	// 检查用户是否已存在
	if err := data.Db.Where("name = ?", username).First(&user).Error; err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exists"},
		})
	} else {
	// 创建用户对象
		newUser := User{
			Name:     username,
			Password: password,
		}
		result := data.Db.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User save failed"},
			})
		} else {
			data.Db.Where("name = ?", username).First(&user)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    username,
			})
		}
	}
}

//var user User

// Login 函数用于处理用户登录的请求。
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//token := username + password

	// 查询用户信息
	var user User
	if err := data.Db.Where("name = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist or invalid credentials"},
		})
		return
	}

	// 验证密码
	if user.Password != password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Invalid credentials"},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    user.Name,
	})
}

// UserInfo 函数用于处理获取用户信息的请求。
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	//id := c.Query("user_id")

  	paasURL := os.Getenv("paas_url")

	if paasURL == "" {
		fmt.Println("环境变量 paas_url 未设置")
		return
	}
	var user User 
	if err := data.Db.Where("name = ?", token).First(&user).Error; err == nil {
    user.Avatar = "https://"+paasURL+"/static/"+user.Avatar
    user.BackgroundImage = "https://"+paasURL+"/static/"+user.BackgroundImage
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

