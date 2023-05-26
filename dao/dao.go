package dao

import (
	"go_blog/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Manager 定义接口，调用对数据库的操作的时候就变成直接调用接口
type Manager interface {
	//用户操作
	// Register 注册
	Register(user *model.User)
	// Login 登录
	Login(username string) model.User
	//修改密码
	ChangePassword(username, newpassword string)

	//  博客操作
	//AddBlog 添加博客
	AddBlog(blog *model.Blog)
	// GetAllBlog 查询所有博客
	GetAllBlog() []model.Blog
	//根据用户名找博客
	GetBlogByUsername(username string) []model.Blog
	// GetBlog 查询单个博客
	GetBlog(pid int) model.Blog
}

//封装数据库的db，将db变成manager，对数据库操作的db换成了manager
type manager struct {
	db *gorm.DB
}

// Mgr 提供外部调用接口
var Mgr Manager

func init() {
	dsn := "root:123456@tcp(127.0.0.1)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	//实例化manager，将Mgr关联到数据库，让Mgr可以用了，相当于将Mgr变成数据库对象，可以让Mgr对数据库进行操作
	Mgr = &manager{db: db}
	//AutoMigrate作用：让数据库自动创建表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&model.Blog{})
	if err != nil {
		return
	}
}

// Register 相当于绑定到数据库方法创建用户
func (mgr *manager) Register(user *model.User) {
	mgr.db.Create(user)
}

func (mgr *manager) Login(username string) model.User {
	var user model.User
	mgr.db.Where("username=?", username).First(&user)
	return user
}

func (mgr *manager) ChangePassword(username, newpassword string) {
	var user model.User
	mgr.db.Model(&user).Where("username=?", username).Update("password", newpassword)
}

func (mgr *manager) AddBlog(blog *model.Blog) {
	mgr.db.Create(blog)
}

func (mgr *manager) GetAllBlog() []model.Blog {
	var blog []model.Blog
	mgr.db.Find(&blog)
	return blog
}

func (mgr *manager) GetBlog(pid int) model.Blog {
	var blog model.Blog
	//mgr.db.Where("id=?",pid).First(&blog)
	mgr.db.First(&blog, pid)
	return blog
}

func (mgr *manager) GetBlogByUsername(username string) []model.Blog {
	var blogs []model.Blog
	mgr.db.Where("username=?", username).Find(&blogs)
	return blogs
}
