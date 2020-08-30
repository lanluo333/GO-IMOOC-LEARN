package middleware

import (
	"github.com/kataras/iris"
)

func AuthConProduct(ctx iris.Context)  {
	uid := ctx.GetCookie("uid")
	// 暂时测试写死
	if uid == "" {
		ctx.Application().Logger().Debug("请先登录再操作")
		//ctx.Redirect("/user/login")
		//return
	}

	ctx.Application().Logger().Debug("登录成功")
	ctx.Next()
}



