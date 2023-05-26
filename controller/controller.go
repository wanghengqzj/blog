package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"go_blog/dao"
	"go_blog/model"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
)

//操作用户

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")

	user := model.User{
		Model:    gorm.Model{},
		Username: username,
		Password: password,
	}

	if repassword != password {
		fmt.Println("两次密码不一致！")
		c.HTML(http.StatusOK, "register.html", "两次密码输入不一致！")
	} else {
		if dao.Mgr.Login(username).Username != "" {
			fmt.Println("该用户名已经被注册！")
			c.HTML(http.StatusOK, "register.html", "该用户名已经被注册！")
		} else {
			fmt.Println("注册成功！")
			dao.Mgr.Register(&user)
			c.Redirect(http.StatusMovedPermanently, "/")
		}
	}
}

func GoRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GoLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println(username)
	user := dao.Mgr.Login(username)

	if user.Username == "" {
		fmt.Println("用户名不存在！")
		c.HTML(http.StatusOK, "login.html", "用户名不存在")
	} else {
		if user.Password != password {
			fmt.Println("密码错误！")
			c.HTML(http.StatusOK, "login.html", "密码错误")
		} else {
			fmt.Println("登录成功！")
			fmt.Println(user.Username)
			cookie, _ := c.Cookie("token")
			if cookie == "" {
				cookie = "NotSet"
				fmt.Println(cookie)
				c.SetCookie("token", user.Username, 3600, "", "localhost", false, true)
			}
			c.Redirect(http.StatusMovedPermanently, "/user_index")
		}
	}
}

func GoUser(c *gin.Context) {
	fmt.Println("---------------")
	cookie, err := c.Cookie("token")
	if err != nil {
		fmt.Println("从cookie中获取用户名失败:", err)
	}
	username := cookie
	fmt.Println("cookie value:", cookie)
	c.HTML(http.StatusOK, "userIndex.html", gin.H{
		"username": username,
	})
}

func Logout(c *gin.Context) {
	//清除cookie
	fmt.Println("用户退出...")
	username, err := c.Cookie("token")
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(username)
	if username != "" {
		fmt.Println("清除cookie...")
		c.SetCookie("token", "", 0, "", "", false, true)
	} else {
		fmt.Println("此时cookie为空...")
	}
	c.Redirect(http.StatusMovedPermanently, "/")
}

func GoUserInfo(c *gin.Context) {
	c.HTML(http.StatusOK, "userInfo.html", nil)
}

func ChangePassword(c *gin.Context) {
	username, err := c.Cookie("token")
	if err != nil {
		fmt.Println("从cookie中获取用户名失败:", err)
	}
	newpassword := c.PostForm("newpassword")
	oldpassword := c.PostForm("password")

	if oldpassword != dao.Mgr.Login(username).Password {
		fmt.Println("原密码错误！")
		c.HTML(http.StatusOK, "userInfo.html", gin.H{
			"msg": "原密码错误，请重试！",
		})
	} else {
		dao.Mgr.ChangePassword(username, newpassword)
		c.Redirect(http.StatusMovedPermanently, "/user_index")
	}
}

//操作博客

func GetBlogIndex(c *gin.Context) {
	blogs := dao.Mgr.GetAllBlog()
	fmt.Println("blogs.....")
	c.HTML(http.StatusOK, "blogIndex.html", blogs)
}

func AddBlog(c *gin.Context) {
	username, err := c.Cookie("token")
	if err != nil {
		fmt.Println("从cookie中取用户名出错！", err)
	}
	title := c.PostForm("title")
	content := c.PostForm("content")
	tag := c.PostForm("tag")
	blog := &model.Blog{
		Model:    gorm.Model{},
		Username: username,
		Title:    title,
		Content:  content,
		Tag:      tag,
	}
	dao.Mgr.AddBlog(blog)

	c.Redirect(http.StatusMovedPermanently, "/blog_by_user")
}

//跳转到添加博客
func GoAddBlog(c *gin.Context) {
	username, err := c.Cookie("token")
	if err != nil {
		fmt.Println("从cookie中取用户名出错！", err)
	}
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"username": username,
	})
}

func BlogDetail(c *gin.Context) {
	s := c.Query("pid")
	pid, _ := strconv.Atoi(s)
	p := dao.Mgr.GetBlog(pid)

	content := blackfriday.Run([]byte(p.Content))

	c.HTML(200, "blogDetail.html", gin.H{
		"Title":   p.Title,
		"Content": template.HTML(content),
	})
}

func GetBlogByUsername(c *gin.Context) {
	username, err := c.Cookie("token")
	if err != nil {
		fmt.Println("从cookie中取用户名出错！", err)
	}
	blogs := dao.Mgr.GetBlogByUsername(username)
	fmt.Println("blogs.....")
	fmt.Println(username)
	fmt.Println(blogs)
	c.HTML(http.StatusOK, "blogByUser.html", blogs)
}
