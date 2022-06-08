package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response        //响应
	UserList []User `json:"user_list"` //用户列表
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")           //查询鉴权
	actionstr := c.Query("action_type") //查询操作行为
	action, err := strconv.Atoi(actionstr)
	if err != nil {
		fmt.Println(err.Error() + "strconv.Atoi(actionstr)")
	}
	touseridstr := c.Query("to_user_id") //查询被操作用户id
	touserid, errs := strconv.Atoi(touseridstr)
	if errs != nil {
		fmt.Println(errs.Error() + "strconv.Atoi(touseridstr)")
	}
	user, _ := QueryToken(token) //根据鉴权寻找用户
	if user.Id > 0 {             //exist 存在的意思 ，判断用户是否存在
		if action == 1 {
			if touserid > 0 && user.Id != int64(touserid) { //不能关注不存在的用户和自己
				SaveRelation(int(user.Id), touserid)
			}
		} else if action == 2 && user.Id != int64(touserid) { //不能取关不存在的用户和自己
			if touserid > 0 {
				DeleteRelation(int(user.Id), touserid)
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0}) //存在，返回正常状态
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"}) //不存在，显示状态码0，并返回StatusMsg，即状态信息：不存在
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) { //关注列表
	token := c.Query("token")
	user, _ := QueryToken(token)
	users, err := GetFollowList(user.Id)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0, //正常状态
		},
		UserList: users, //视频列表
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) { //粉丝列表
	token := c.Query("token")
	user, _ := QueryToken(token)
	users, err := GetFollowerList(user.Id)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0, //正常状态
		},
		UserList: users, //用户列表
	})
}

func SaveRelation(userid int, touserid int) {
	sqlstr := "insert into relation(userid,touserid) values(?,?)"
	_, err := Db.Exec(sqlstr, userid, touserid)
	if err != nil {
		fmt.Println(err.Error() + "SaveRelation")
	}
	sqlstrs := "update users set followercount = followercount +1  where id=? "
	_, errs := Db.Exec(sqlstrs, touserid)
	if errs != nil {
		fmt.Println(err.Error() + "SaveRelation2")
	}

	sqlstrss := "update users set followcount=followcount +1 where id=? "
	_, errss := Db.Exec(sqlstrss, userid)
	if errss != nil {
		fmt.Println(err.Error() + "SaveRelation3")
	}
}
func DeleteRelation(userid int, touserid int) {
	sqlstr := "delete from relation where userid= ? and touserid= ?"
	_, err := Db.Exec(sqlstr, userid, touserid)
	if err != nil {
		fmt.Println(err.Error() + "DeleteRelation")
	}
	sqlstrs := "update users set followercount = followercount -1 where id=? and followercount > 0"
	_, errs := Db.Exec(sqlstrs, touserid)
	if errs != nil {
		fmt.Println(err.Error() + "DeleteRelation2")
	}
	sqlstrss := "update users set followcount = followcount -1 where id=? and followcount > 0"
	_, errss := Db.Exec(sqlstrss, userid)
	if errss != nil {
		fmt.Println(err.Error() + "SaveRelation3")
	}
}

func GetFollowList(userid int64) ([]User, error) {
	var Users []User
	var ToUserids []int
	var touserid int
	var user User
	{ //获取关注列表id集合
		sqlstr := "select touserid from relation where userid= ?  and userid > 0" //获取自己关注的用户id集合
		rows, err := Db.Query(sqlstr, userid)
		if err != nil {
			fmt.Println(err.Error() + "GetUserList.Db")
			return nil, err
		}
		for rows.Next() {
			err := rows.Scan(&touserid)
			if err != nil {
				return nil, err
			}
			ToUserids = append(ToUserids, touserid)
		}
	}
	{ //根据id集合获取用户信息集合
		for touser := range ToUserids {
			sqlstr := "select id,name,followcount,followercount,token from users where id=?"
			row := Db.QueryRow(sqlstr, touser)
			row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.Token)
			if err != nil {
				fmt.Println(err.Error() + "GetUserList.Db")
				return nil, err
			}
			var ra Relation
			sqlstrs := "select id from relation where touserid= ? userid =?"
			rowss := Db.QueryRow(sqlstrs, touser, userid)
			rowss.Scan(&ra.Id)
			if ra.Id > 0 {
				user.IsFollow = true
			}
			Users = append(Users, user)
		}
	}
	return Users, nil
}
func GetFollowerList(userid int64) ([]User, error) {
	var Users []User
	var Userids []int
	var userids int
	var user User
	{ //获取粉丝列表id集合
		sqlstr := "select userid from relation where touserid= ? and  touserid > 0 " //获取关注自己的用户的id集合
		rows, err := Db.Query(sqlstr, userid)
		if err != nil {
			fmt.Println(err.Error() + "GetUserList.Db")
			return nil, err
		}
		for rows.Next() {
			err := rows.Scan(&userids)
			if err != nil {
				return nil, err
			}
			Userids = append(Userids, userids)
		}
	}
	{ //根据id集合获取用户信息集合
		for userids := range Userids {
			sqlstr := "select id,name,followcount,followercount,token from users where id=?"
			row := Db.QueryRow(sqlstr, userids)
			row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.Token)
			if err != nil {
				fmt.Println(err.Error() + "GetUserList.Db")
				return nil, err
			}
			Users = append(Users, user)
		}
	}
	return Users, nil
}
