package controller


// Response 结构体表示API接口返回的通用响应数据。
type Response struct {
	StatusCode int32  `json:"status_code"`          // 响应状态码
	StatusMsg  string `json:"status_msg,omitempty"` // 响应状态消息，可选字段
}

// Video 结构体表示视频信息。
type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author" gorm:"foreignKey:AuthorId"`   //视频作者
	AuthorId      int64  `json:"aid"` 
	PlayUrl       string `json:"play_url,omitempty"`         //视频url
	CoverUrl      string `json:"cover_url,omitempty"`        //封面url
	FavoriteCount int64  `json:"favorite_count,omitempty"`   //点赞数
	CommentCount  int64  `json:"comment_count,omitempty"`    //评论数
	IsFavorite    bool   `json:"is_favorite,omitempty"`     //是否被点赞
	Title         string `json:"title"`    //标题
}


// Comment 结构体表示评论信息。
type Comment struct {
	Id         int64  `json:"id,omitempty"`          // 评论ID
	User       User   `json:"user"`                  // 评论用户信息
	Content    string `json:"content,omitempty"`     // 评论内容
	CreateDate string `json:"create_date,omitempty"` // 评论创建日期
}



// User 结构体表示用户信息。
type User struct {
	Id              int64  `json:"id,omitempty"`               // 用户ID
	Name            string `json:"name,omitempty"`             // 用户名
	Password        string `json:"password,omitempty"`         //用户密码
	FollowCount     int64  `json:"follow_count,omitempty"`     // 关注数
	FollowerCount   int64  `json:"follower_count,omitempty"`   // 粉丝数
	IsFollow        bool   `json:"is_follow,omitempty"`        // 是否已关注该用户
	Avatar          string `json:"avatar,omitempty" gorm:"default:qq_pic_merged_1633959190718.jpg"`   //用户头像
	BackgroundImage string `json:"background_image,omitempty" gorm:"default:QQpic20211019200841.jpg"` //背景图片
	Signature       string `json:"signature,omitempty" gorm:"default:my signature hhh"`        //签名
	TotalFavorited  string `json:"total_favorited,omitempty"`  
	WorkCount       int64  `json:"work_count,omitempty"`       //作品数
	FavoriteCount   int64  `json:"favorite_count,omitempty"`   //喜欢数
}

// VideoListItem 自定义结构体用于映射视频列表项
type VideoListItem struct {
	ID            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

// FeedResponse 自定义结构体用于接口响应
type FeedResponse struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	NextTime   int64           `json:"next_time"`
	VideoList  []VideoListItem `json:"video_list"`
}

// Message 结构体表示聊天消息。
type Message struct {
	Id         int64  `json:"id,omitempty"`          // 消息ID
	Content    string `json:"content,omitempty"`     // 消息内容
	CreateTime string `json:"create_time,omitempty"` // 消息创建时间
}

// MessageSendEvent 结构体表示发送聊天消息的事件。
type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`     // 发送消息的用户ID
	ToUserId   int64  `json:"to_user_id,omitempty"`  // 接收消息的用户ID
	MsgContent string `json:"msg_content,omitempty"` // 消息内容
}

// MessagePushEvent 结构体表示推送聊天消息的事件。
type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`     // 消息发送者的用户ID
	MsgContent string `json:"msg_content,omitempty"` // 消息内容
}