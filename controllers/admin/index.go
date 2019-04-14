package admin

import (
	"os"
	"runtime"
	"server/models"

	"github.com/astaxie/beego"
)

type IndexController struct {
	baseController
}

func (this *IndexController) Index() {
	beego.Debug("进入后台")
	this.Data["version"] = beego.AppConfig.String("AppVer")
	this.Data["adminid"] = this.userid
	this.Data["adminname"] = this.username
	this.TplName = this.moduleName + "/index/index.html"
}

func (this *IndexController) VueIndex() {
	beego.Debug("进入后台")
	param := make(map[string]interface{})
	param["state"] = 1
	param["hostname"], _ = os.Hostname()
	param["version"] = beego.AppConfig.String("AppVer")
	param["gover"] = runtime.Version()
	param["os"] = runtime.GOOS
	param["cpunum"] = runtime.NumCPU()
	param["arch"] = runtime.GOARCH
	param["ip"] = this.getClientIp()
	param["uid"] = this.userid
	param["userName"] = this.username

	beego.Debug(this.getClientIp())

	this.Data["json"] = param
	this.ServeJSON()
}

func (this *IndexController) Main() {
	beego.Debug("登陆成功")
	this.Data["hostname"], _ = os.Hostname()
	this.Data["version"] = beego.AppConfig.String("AppVer")
	this.Data["gover"] = runtime.Version()
	this.Data["os"] = runtime.GOOS
	this.Data["cpunum"] = runtime.NumCPU()
	this.Data["arch"] = runtime.GOARCH

	this.Data["postnum"], _ = new(models.Post).Query().Count()
	this.Data["tagnum"], _ = new(models.Tag).Query().Count()
	this.Data["usernum"], _ = new(models.User).Query().Count()

	this.display()
}
