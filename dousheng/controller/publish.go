package controller

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token") //从表单获取鉴权
	title := c.PostForm("title")
	user, err := QueryToken(token) //获取用户信息
	if err != nil {
		fmt.Println(err.Error())
	}
	if user.Id <= 0 { //如果用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"}) //不存在，显示状态码1，并返回StatusMsg，即状态信息
		return
	}
	data, err := c.FormFile("data") //从表单获取文件
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,           //异常状态码
			StatusMsg:  err.Error(), //返回错误信息
		})
		return
	}
	Coverurl := "http://192.168.0.103:8080/static/public/douyin.jpeg"
	fmt.Println(title, "a")
	filename := filepath.Base(data.Filename)                       //获取文件名
	finalName := fmt.Sprintf("%v%v%v", user.Name, title, filename) //获取文件标准格式
	saveFile := fmt.Sprintf("%v%v", "public/", finalName)          //配置保存文件在public下json字段
	urlpath := "http://192.168.0.103:8080/static/"
	Playurl := fmt.Sprintf("%v%v", urlpath, finalName)
	SaveByAuthor(user.Id, Playurl, Coverurl, title)
	if err := c.SaveUploadedFile(data, saveFile); err != nil { //保存文件 ，如果保存错误
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,           //异常状态码
			StatusMsg:  err.Error(), //返回错误信息0
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,                                //正常状态码
		StatusMsg:  title + " uploaded successfully", //返回成功信息
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) { //个人界面作品列表
	fmt.Println("来自PublishList")
	token := c.Query("token")
	user, err := QueryToken(token) //获取用户信息
	if err != nil {
		fmt.Println(err.Error())
	}
	Videos, err := GetVideoByPublish(user.Id)
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

func GetVideoByPublish(authorid int64) ([]Video, error) {
	var videos []Video
	var video Video
	sqlstr := "select videoid,userid,coverurl,favoritecount,commentcount,title from video where userid =?"
	rows, err := Db.Query(sqlstr, authorid)
	if err != nil {
		fmt.Println(err.Error(), "rows")
	}
	for rows.Next() {
		errs := rows.Scan(&video.Id, &video.Author.Id, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.Title)
		if errs != nil {
			fmt.Println(errs.Error())
			return videos, errs
		}
		video.PlayUrl = ""
		videos = append(videos, video)
	}
	return videos, nil
}

func SaveByAuthor(authorid int64, playurl string, Coverurl string, title string) {
	sqlstr := "insert into video(userid,playurl,coverurl,favoritecount,commentcount,title) values(?,?,?,?,?,?)"
	_, err := Db.Exec(sqlstr, authorid, playurl, Coverurl, 0, 0, 0, title)
	if err != nil {
		fmt.Println(err.Error() + " SaveByAuthor")
	}
}
