package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/douyin")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("数据库链接成功")
}

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyid
/*var usersLoginInfo = map[string]User{ //【鉴权】（token）：用户（User）
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}
*/
//var userIdSequence = int64(1) //用户ID序列号

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"` //用户ID
	Token  string `json:"token"`             //用户鉴权token
}

type UserResponse struct {
	Response      //响应行
	User     User `json:"user"` //用户信息
}

//注册账户
func Register(c *gin.Context) {
	username := c.Query("username") //获取登录账户
	password := c.Query("password") //获取登录密码

	token := username + password //计算鉴权
	user, _ := QueryToken(token)
	fmt.Println("用户是", user)
	if user.Id > 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"}, //存在，返回StutasMsg：用户已存在
		})
	} else { //若用户不存在，则注册
		SavaUser(username, token)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0}, //响应行
			UserId:   user.Id,                 //返回UserID
			Token:    token,                   //返回鉴权
		})
	}
}

//登录账户
func Login(c *gin.Context) {
	username := c.Query("username") //获取登录账户
	password := c.Query("password") //获取登录密码

	token := username + password //计算鉴权
	user, _ := QueryToken(token)
	fmt.Println(user.Id)
	if user.Id > 0 { //exist 存在的意思 ，判断用户是否存在
		fmt.Println("登录成功！")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0}, //响应头,0表示正确,即用户存在
			UserId:   user.Id,                 //返回用户ID
			Token:    token,                   //返回鉴权
		})
	} else {
		fmt.Println("登录失败！")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"}, //不存在，显示状态码0，并返回StatusMsg，即状态信息
		})
	}
}

//返回用户信息
func UserInfo(c *gin.Context) {
	token := c.Query("token") //在URL中查询鉴权,本demo的token为username+password
	user, _ := QueryToken(token)
	if user.Id > 0 { //exist 存在的意思 ，判断用户是否存在
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0}, //存在，状态码为0，即正常
			User:     user,                    //返回User信息
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"}, //不存在，显示状态码0，并返回StatusMsg，即状态信息
		})
	}
}

//根据token查询用户
func QueryToken(token string) (User, error) {
	var user User
	sqlstr := "select id,name from users where token=?"
	row := Db.QueryRow(sqlstr, token)
	row.Scan(&user.Id, &user.Name)
	return user, nil
}

//保存用户信息
func SavaUser(username string, token string) error {
	sqlStr := "insert into users(name,followcount,followercount,follow,token) values(?,?,?,?,?)"
	_, err := Db.Exec(sqlStr, username, 0, 0, 0, token)
	if err != nil {
		return err
	}
	return nil
}
