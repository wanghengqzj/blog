package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_blog/controller"
	"html/template"
	"net/http"
)

func MyHandler(c *gin.Context) {
	fmt.Println("MyHandler....")
}

func MyHandlerB() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()
		meth := context.Request.Method
		fmt.Println("MyHandlerB....", path, "  ", meth)
	}
}

//go:embed templates/* assets/*
var f embed.FS

func main() {
	e := gin.Default()
	e.Use(MyHandler, MyHandlerB())

	html := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	e.SetHTMLTemplate(html)
	//e.Any("/static/*filepath", func(context *gin.Context) {
	//	staticServer := http.FileServer(http.FS(f))
	//	staticServer.ServeHTTP(context.Writer, context.Request)
	//})

	e.NoRoute(gin.WrapH(http.FileServer(http.FS(f))))

	//e.LoadHTMLGlob("templates/*")
	//e.Static("/assets", "./assets")

	//链接都是get请求，从网页链接的get请求跳转到post请求
	e.POST("/register", controller.Register)
	e.GET("/register", controller.GoRegister)
	e.GET("/login", controller.GoLogin)
	e.POST("/login", controller.Login)
	e.GET("/user_index", controller.GoUser)
	e.POST("/logout", controller.Logout)
	e.GET("/user_info", controller.GoUserInfo)
	e.POST("/user_info", controller.ChangePassword)

	//博客列表
	e.GET("/blog_index", controller.GetBlogIndex)
	//根据用户名返回博客列表
	e.GET("/blog_by_user", controller.GetBlogByUsername)
	//跳转到添加博客页面
	e.GET("/blog", controller.GoAddBlog)
	//添加博客
	e.POST("/blog", controller.AddBlog)
	//跳转到详细页面
	e.GET("/blog_detail", controller.BlogDetail)

	e.GET("/", controller.Index)

	err := e.Run()
	if err != nil {
		return
	}
}
