package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response              //响应行
	CommentList []Comment `json:"comment_list,omitempty"` //评论列表
}

type CommentActionResponse struct {
	Response         //响应行
	Comment  Comment `json:"comment,omitempty"` //评论实体
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")            //查询鉴权，即用户登录才能评论
	actionType := c.Query("action_type") //动作类型,1-发布评论，2-删除评论
	user, _ := QueryToken(token)
	if user.Id > 0 { //exist 存在的意思 ，判断用户是否存在
		time := fmt.Sprintf("%v-%v-%v %v:%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day(), time.Now().Hour(), time.Now().Minute())
		if actionType == "1" {
			videoidstr := c.Query("video_id")
			videoid, err := strconv.Atoi(videoidstr)
			if err != nil {
				fmt.Println("v" + err.Error())
			}
			text := c.Query("comment_text") //查询评论内容
			errs := SaveComment(user.Id, text, time, videoid)
			if errs != nil {
				fmt.Println(errs.Error())
			}
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0}, //正常状态码
				Comment: Comment{ //返回评论信息
					Id:         user.Id,
					User:       user,
					Content:    text,
					CreateDate: time,
				}})
			return
		} else if actionType == "2" {
			videoidstr := c.Query("video_id")
			videoid, err := strconv.Atoi(videoidstr)
			if err != nil {
				fmt.Println("v" + err.Error())
			}
			Commentidstr := c.Query("comment_id")
			Commentid, err := strconv.Atoi(Commentidstr)
			if err != nil {
				fmt.Println("c" + err.Error())
			}
			errs := DeleteComment(Commentid, videoid, int(user.Id))
			if errs != nil {
				fmt.Println(errs.Error())
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0}) //正常状态码
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"}) //异常状态码，返回错误信息
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) { //评论列表
	videoidstr := c.Query("video_id")
	videoid, err := strconv.Atoi(videoidstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	Comments, err := GetCommentList(videoid)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0}, //正常状态码
		CommentList: Comments,                //返回评论列表
	})
}

func SaveComment(userid int64, text string, time string, videoid int) error {
	sqlstr := "insert into comment (userid,content,createdate,videoid) values(?,?,?,?)"
	_, err := Db.Exec(sqlstr, userid, text, time, videoid)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlstr2 := "update video set commentcount=commentcount+1  where videoid=?"
	_, errs := Db.Exec(sqlstr2, videoid)
	if errs != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func GetCommentList(videoid int) ([]Comment, error) {
	var Comments []Comment
	var Comment Comment
	sqlstr := "select * from comment where videoid = ?"
	rows, err := Db.Query(sqlstr, videoid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&Comment.Id, &Comment.User.Id, &Comment.Content, &Comment.CreateDate, &videoid)
		if err != nil {
			return Comments, err
		}
		var users User
		sqlstrs := "select * from users where id= ?"
		row := Db.QueryRow(sqlstrs, Comment.User.Id)
		row.Scan(&users.Id, &users.Name, &users.FollowCount, &users.FollowerCount, &users.Token)
		fmt.Println(Comment.User.Id, users.Name)
		Comment.User = users
		Comments = append(Comments, Comment)
	}
	return Comments, nil
}

func DeleteComment(Commentid int, videoid int, userid int) error {
	sqlstr := "delete from comment where id = ? and userid= ?"
	_, err := Db.Exec(sqlstr, Commentid, userid)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlstr2 := "update video set commentcount=commentcount-1  where videoid=? and  commentcount>0"
	_, errs := Db.Exec(sqlstr2, videoid)
	if errs != nil {
		fmt.Println(err.Error())
	}
	return nil
}
