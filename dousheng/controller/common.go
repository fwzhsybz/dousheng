package controller

type Response struct {
	StatusCode int32  `json:"status_code"`          //状态码
	StatusMsg  string `json:"status_msg,omitempty"` //返回的状态消息
}

type Video struct {
	Id            int64  `json:"id,omitempty"`             //视频唯一标识
	Author        User   `json:"author"`                   //作者信息
	PlayUrl       string `json:"play_url,omitempty"`       //视屏播放地址
	CoverUrl      string `json:"cover_url,omitempty"`      //视频封面地址
	FavoriteCount int64  `json:"favorite_count,omitempty"` //视频点赞总数
	CommentCount  int64  `json:"comment_count,omitempty"`  //视频评论总数
	IsFavorite    bool   `json:"is_favorite,omitempty"`    //是否点赞
	Title         string `json:"title"`                    //标题
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`          //评论ID
	User       User   `json:"user"`                  //评论用户的信息
	Content    string `json:"content,omitempty"`     //评论内容
	CreateDate string `json:"create_date,omitempty"` //评论日期
}

type User struct {
	Id            int64  `json:"id,omitempty"`             //用户ID
	Name          string `json:"name,omitempty"`           //用户名称
	FollowCount   int64  `json:"follow_count,omitempty"`   //关注数
	FollowerCount int64  `json:"follower_count,omitempty"` //粉丝数
	IsFollow      bool   `json:"is_follow,omitempty"`      //是否关注
	Token         string `json:"token"`
}

type Favorite struct {
	Id      int
	Userid  int
	Videoid int
}

type Relation struct {
	Id       int
	Userid   int
	Touserid int
}
