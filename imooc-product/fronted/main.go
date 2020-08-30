package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-shop/common"
	"imooc-shop/fronted/middleware"
	"imooc-shop/fronted/web/controllers"
	"imooc-shop/rabbitmq"
	"imooc-shop/repositories"
	"imooc-shop/services"
)

func main()  {
	// 1. 创建iris实例
	app := iris.New()
	// 2. 设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	// 3. 设置模板目标
	app.StaticWeb("/public", "./fronted/web/public")
	// 访问生成好的html静态文件
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")

	// 4. 注册模板
	template := iris.HTML("./fronted/web/views", ".html").
		Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println("db error")
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	rabbitMq := rabbitmq.NewRabbitMQSimple("imoocProduct")

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(*product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(*order)
	proProduct := app.Party("/product")
	proProduct.Use(middleware.AuthConProduct)
	pro := mvc.New(proProduct)
	pro.Register(productService, orderService, ctx, rabbitMq)
	pro.Handle(new(controllers.ProductController))

	app.Run(
			iris.Addr(":8082"),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
		)

}


