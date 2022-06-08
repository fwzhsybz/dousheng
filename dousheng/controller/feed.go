package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"` //视频列表
	NextTime  int64   `json:"next_time,omitempty"`  //本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	fmt.Println("来自Feed")
	token := c.Query("token")
	user, _ := QueryToken(token)
	Videos, err := GetVideo(int(user.Id))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("普通用户获得的视频")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0}, //正常状态
		VideoList: Videos,                  //视频列表
		NextTime:  time.Now().Unix(),       //初始化为当前时间
	})
}

func GetVideo(userid int) ([]Video, error) {
	var videos []Video
	var video Video
	sqlstr := "select * from video where videoid >?"
	rows, err := Db.Query(sqlstr, 0)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&video.Id, &video.Author.Id, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.Title)
		if err != nil {
			return videos, err
		}
		var favorite Favorite
		sqlstr := "select * from favorite where videoid= ? and userid= ?"
		rowss := Db.QueryRow(sqlstr, video.Id, userid)
		rowss.Scan(&favorite.Id, &favorite.Userid, &favorite.Videoid)
		if favorite.Id > 0 {
			video.IsFavorite = true
		} else {
			video.IsFavorite = false
		}
		var ra Relation
		sqlstrss := "select id from relation where touserid= ? userid =?"
		rowsss := Db.QueryRow(sqlstrss, video.Author.Id, userid)
		rowsss.Scan(&ra.Id)
		if ra.Id > 0 {
			video.Author.IsFollow = true
		}
		var user User
		sqlstrs := "select * from users where id= ?"
		row := Db.QueryRow(sqlstrs, video.Author.Id)
		row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.Token)
		video.Author = user

		videos = append(videos, video)
	}
	return videos, nil
}
