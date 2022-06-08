package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid //简单demo功能不全
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")           //查询鉴权
	actionstr := c.Query("action_type") //查看操作
	action, err := strconv.Atoi(actionstr)
	if err != nil {
		fmt.Println("a" + err.Error())
	}
	videoidstr := c.Query("video_id")
	videoid, errs := strconv.Atoi(videoidstr)
	if errs != nil {
		fmt.Println("v" + err.Error())
	}
	user, _ := QueryToken(token)
	if user.Id > 0 { //exist 存在的意思 ，判断用户是否存在
		if action == 1 {
			ActionForFavorite(1, videoid, int(user.Id))
		} else if action == 2 {
			ActionForFavorite(0, videoid, int(user.Id))
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0}) //存在，状态码为0，即正常
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"}) //不存在，显示状态码0，并返回StatusMsg，即状态信息
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) { //点赞列表
	var Videos []Video
	token := c.Query("token")
	user, _ := QueryToken(token)
	Videos, err := GetVideoByFavorite(int64(user.Id))
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0, //正常状态码
		},
		VideoList: Videos, //视频列表
	})
}

func ActionForFavorite(action int, videoid int, userid int) {
	sqlstrs := "select id from favorite where userid=? and videoid =?"
	row := Db.QueryRow(sqlstrs, userid, videoid)
	var fa Favorite
	row.Scan(&fa.Id)
	if action == 1 && fa.Id <= 0 {
		sqlstrs := "insert into favorite(userid,videoid) values(?,?)"
		_, errs := Db.Exec(sqlstrs, userid, videoid)
		if errs != nil {
			fmt.Println(errs.Error() + "insert favorite")
		}
		sqlstr := "update video set favoritecount =  favoritecount +1 where videoid=? "
		_, err := Db.Exec(sqlstr, videoid)
		if err != nil {
			fmt.Println(err.Error() + "a1")
		}
	} else if action == 0 && fa.Id > 0 {
		sqlstrs := "delete from favorite where userid=? and videoid =?"
		_, errs := Db.Exec(sqlstrs, userid, videoid)
		if errs != nil {
			fmt.Println(errs.Error() + "insert favorite")
		}
		sqlstr := "update video set favoritecount =  favoritecount -1 where videoid=?"
		_, err := Db.Exec(sqlstr, videoid)
		if err != nil {
			fmt.Println(err.Error() + "a2")
		}
	}
}
func GetVideoByFavorite(userid int64) ([]Video, error) {
	var videos []Video
	var video Video
	sqlstr := "select videoid,userid,coverurl,favoritecount,commentcount,title from video where videoid >?" //获取信息从video库里
	rows, err := Db.Query(sqlstr, 0)
	if err != nil {
		fmt.Println(err.Error(), "rows")
	}
	for rows.Next() {
		errs := rows.Scan(&video.Id, &video.Author.Id, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.Title)
		if errs != nil {
			fmt.Println(errs.Error())
			return videos, errs
		}
		var favorite Favorite
		sqlstr := "select * from favorite where videoid= ? and userid= ?"
		rowss := Db.QueryRow(sqlstr, video.Id, userid)
		rowss.Scan(&favorite.Id, &favorite.Userid, &favorite.Videoid)
		video.PlayUrl = ""
		if favorite.Id > 0 {
			videos = append(videos, video)
		}
	}
	return videos, nil
}
