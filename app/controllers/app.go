package controllers

import (
	"github.com/robfig/revel"
	avatar "avatar/app/engine"
)

const (
	EMAIL_STR = "@hotmail.com"
)

type App struct {
	*revel.Controller
}

func init() {
	revel.OnAppStart(func() {
		avatar.StartEngine()
	})
}
/*
func (c App) Stop() revel.Result{
	avatar.ControlAvatar("stop")
	return c.Render()
}
*/
func (c App) Index() revel.Result {
	return c.Render()
}


func (c App) Welcome() revel.Result {
	//revel.INFO.Println(c.Session["user"])
	return c.Render()
}

//用户登陆
func (c App) Login(userName, passwd string) revel.Result {
	c.Validation.Email(userName+EMAIL_STR).Message("Your email address is invalid")
	c.Validation.Required(passwd).Message("Please enter user passwd")
	if c.Validation.HasErrors() {
		c.FlashParams()
		c.Validation.Keep()
		return c.Redirect(App.Index)
	}

	c.Session["user"] = userName
	//c.Flash.Success("login successful")
	//return c.Redirect("/welcome/%s" userName)
	return c.Redirect("/welcome")
}

