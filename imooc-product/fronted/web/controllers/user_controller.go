package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"imooc-shop/encrypt"
	"imooc-shop/services"
	"strconv"
)

type UserController struct {
	Ctx		iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.Result  {
	return mvc.View{
		Name:"user/register.html",
	}
}

func (c *UserController) PostRegister()  {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user := &datamodels.User{
		NickName:    nickName,
		UserName:    userName,
		HasPassword: password,
	}

	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}

	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name:"user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response  {
	// 1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	// 2. 验证账号密码正确性
	user, isOK := c.Service.IsPwdSuccess(userName, password)
	if !isOK {
		return mvc.Response{
			Path:"/user/login",
		}
	}

	// 写入用户id到cookie中
	//c.Ctx.SetCookie(&http.Cookie{Name:"uid", Value:"1"})
	common.GlobalCookie(c.Ctx, "uid", "1")
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//c.Session.Set("uid", strconv.FormatInt(1, 10))
	// 写入用户浏览器
	common.GlobalCookie(c.Ctx, "sign", uidString)

	return mvc.Response{
		Path:"/product/detail",
	}
}

