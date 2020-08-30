package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-shop/backend/web/controllers"
	"imooc-shop/common"
	"imooc-shop/repositories"
	"imooc-shop/services"
	"log"
)

func main()  {
	// 1. 创建iris实例
	app := iris.New()
	// 2. 设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	// 3. 设置模板目标
	app.StaticWeb("/assets", "./backend/web/assets")

	// 4. 注册模板
	template := iris.HTML("./backend/web/views", ".html").
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
		fmt.Println(err)
		log.Println(err)
	}

	// 创建上下文
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// 5. 注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(*productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(*orderRepository)
	order := mvc.New(app.Party("/order"))
	order.Register(orderService)
	order.Handle(new(controllers.OrderController))

	// 6.启动服务
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed),iris.WithOptimizations)
}


