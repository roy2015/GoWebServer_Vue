package main

import (
	"GoWebServer_Vue/controllers/admin"
	"GoWebServer_Vue/controllers/blog"
	"GoWebServer_Vue/models"
	"GoWebServer_Vue/util"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	util.Factory.Set("cache", func() (interface{}, error) {
		mc := util.NewLruCache(1000)
		return mc, nil
	})

	models.Init()
}

func main() {
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	//前台路由
	beego.Router("/", &blog.MainController{}, "*:Index")
	beego.Router("/page/:page:int", &blog.MainController{}, "*:Index")
	beego.Router("/article/:id:int", &blog.MainController{}, "*:Show")      //ID访问
	beego.Router("/article/:urlname(.+)", &blog.MainController{}, "*:Show") //别名访问
	beego.Router("/archives", &blog.MainController{}, "*:Archives")
	beego.Router("/archives/page/:page:int", &blog.MainController{}, "*:Archives")
	beego.Router("/category/:name(.+?)", &blog.MainController{}, "*:Category")
	beego.Router("/category/:name(.+?)/page/:page:int", &blog.MainController{}, "*:Category")
	beego.Router("/:urlname(.+)", &blog.MainController{}, "*:Show") //别名访问

	//后台路由
	beego.Router("/vueAdmin", &admin.IndexController{}, "*:VueIndex")
	beego.Router("/admin", &admin.IndexController{}, "*:Index")
	beego.Router("/admin/main", &admin.IndexController{}, "*:Main")

	beego.Router("/admin/ShowLogin", &admin.AccountController{}, "*:ShowLogin")
	beego.Router("/admin/login", &admin.AccountController{}, "*:Login")
	beego.Router("/admin/vueLogin", &admin.AccountController{}, "*:VueLogin")
	beego.Router("/admin/logout", &admin.AccountController{}, "*:Logout")
	beego.Router("/admin/account/profile", &admin.AccountController{}, "*:Profile")

	//系统管理
	beego.Router("/admin/system/setting", &admin.SystemController{}, "*:Setting")
	//内容管理
	beego.Router("/admin/article/list", &admin.ArticleController{}, "*:List")
	beego.Router("/admin/article/add", &admin.ArticleController{}, "*:Add")
	beego.Router("/admin/article/edit", &admin.ArticleController{}, "*:Edit")
	beego.Router("/admin/article/save", &admin.ArticleController{}, "post:Save")
	beego.Router("/admin/article/delete", &admin.ArticleController{}, "*:Delete")
	beego.Router("/admin/article/batch", &admin.ArticleController{}, "*:Batch")
	beego.Router("/admin/article/upload", &admin.ArticleController{}, "*:Upload")
	beego.Router("/admin/tag", &admin.TagController{}, "*:Index")
	//用户管理
	beego.Router("/admin/user/vueList", &admin.UserController{}, "*:VueList")
	beego.Router("/admin/user/vueEdit", &admin.UserController{}, "*:VueEdit")
	beego.Router("/admin/user/list", &admin.UserController{}, "*:List")
	beego.Router("/admin/user/add", &admin.UserController{}, "*:Add")
	beego.Router("/admin/user/edit", &admin.UserController{}, "*:Edit")
	beego.Router("/admin/user/delete", &admin.UserController{}, "*:Delete")

	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.Debug("log system test!")
	beego.Debug("Welcom to web back system")
	beego.Run()
}
