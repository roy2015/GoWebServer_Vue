package admin

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lisijie/goblog/controllers/admin/test"
	"github.com/lisijie/goblog/models"
	"github.com/lisijie/goblog/util"
)

type AccountController struct {
	baseController
	test.BaseTest
}

func (this *AccountController) ShowLogin() {
	beego.Debug("进入 AccountController:")
	this.TplName = this.moduleName + "/account/index.html"
}

//登录
func (this *AccountController) Login() {
	if this.userid > 0 {
		this.Redirect("/admin", 302)
	}

	if this.GetString("dosubmit") == "yes" {
		account := strings.TrimSpace(this.GetString("account"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if account != "" && password != "" {
			var user models.User
			user.UserName = account
			beego.Debug("account:", account, "password:", password)
			if user.Read("user_name") != nil || user.Password != util.Md5([]byte(password)) {
				this.Data["errmsg"] = "帐号或密码错误"
			} else if user.Active == 0 {
				this.Data["errmsg"] = "该帐号未激活"
			} else {
				this.Data["errmsg"] = "登陆成功"
				user.LoginCount += 1
				user.LastIp = this.getClientIp()
				user.LastLogin = this.getTime()
				user.Update()
				authkey := util.Md5([]byte(this.getClientIp() + "|" + user.Password))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}
				this.Redirect("/admin", 302)
			}
		}
	}
	this.TplName = this.moduleName + "/account/login.html"
}

//vue方式登录
func (this *AccountController) VueLogin() {
	beego.Debug("this.userid:", this.userid)
	beego.Debug("---- 进入登陆界面 ---")

	if this.userid > 0 {
		//用户存在直接登录
		this.Data["json"] = map[string]interface{}{"success": 0, "state": 1}
		this.ServeJSON()
		return
	}

	if this.GetString("dosubmit") == "yes" {
		account := strings.TrimSpace(this.GetString("account"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if account != "" && password != "" {
			var user models.User
			user.UserName = account
			message := ""
			state := 0
			if user.Read("user_name") != nil || user.Password != util.Md5([]byte(password)) {
				message = "帐号或密码错误"
				state = -1
			} else if user.Active == 0 {
				message = "该帐号未激活"
				state = -1
			} else {
				message = "登陆成功"
				state = 1
				user.LoginCount += 1
				user.LastIp = this.getClientIp()
				user.LastLogin = this.getTime()
				user.Update()
				authkey := util.Md5([]byte(this.getClientIp() + "|" + user.Password))

				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}
				beego.Debug("--登录成功--", user.Id)
			}

			this.Data["json"] = map[string]interface{}{"state": state, "message": message}
			this.ServeJSON()
			return
		}
	}

	//重新登陆
	this.Data["json"] = map[string]interface{}{"state": -1, "message": "请登录"}
	this.ServeJSON()
}

//退出登录
func (this *AccountController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.Redirect("/admin/login", 302)
}

//资料修改
func (this *AccountController) Profile() {
	user := models.User{Id: this.userid}
	if err := user.Read(); err != nil {
		this.showmsg(err.Error())
	}
	if this.isPost() {
		errmsg := make(map[string]string)
		password := strings.TrimSpace(this.GetString("password"))
		newpassword := strings.TrimSpace(this.GetString("newpassword"))
		newpassword2 := strings.TrimSpace(this.GetString("newpassword2"))
		updated := false
		if newpassword != "" {
			if password == "" || util.Md5([]byte(password)) != user.Password {
				errmsg["password"] = "当前密码错误"
			} else if len(newpassword) < 6 {
				errmsg["newpassword"] = "密码长度不能少于6个字符"
			} else if newpassword != newpassword2 {
				errmsg["newpassword2"] = "两次输入的密码不一致"
			}
			if len(errmsg) == 0 {
				user.Password = util.Md5([]byte(newpassword))
				user.Update("password")
				updated = true
			}
		}
		this.Data["updated"] = updated
		this.Data["errmsg"] = errmsg
	}
	this.Data["user"] = user
	this.display()
}
