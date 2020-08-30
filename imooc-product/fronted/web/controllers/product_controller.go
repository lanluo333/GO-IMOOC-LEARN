package controllers

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"html/template"
	"imooc-shop/datamodels"
	"imooc-shop/rabbitmq"
	"imooc-shop/services"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx 			iris.Context
	ProductService	services.ProductService
	OrderService	services.IOrderService
	RabbitMQ 		*rabbitmq.RabbitMQ
	Session 		*sessions.Session
}

var (
	htmlOutPath = "./fronted/web/htmlProductShow/"  // 生成的html保存目录
	templatePath = "./fronted/web/views/template/"  // 静态文件模板目录
)

// 生成静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product)  {
	// 判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}

	// 生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product)

}

func (p *ProductController) GetGenerateHtml()  {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	// 1.获取模板文件地址
	contentsTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	// 2.获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	// 3.获取模板渲染数据
	product, err := p.ProductService.GetProductById(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	// 生成静态文件
	generateStaticHtml(p.Ctx, contentsTmp, fileName, product)
}

// 判断文件是否存在
func exist(fileName string) bool  {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductById(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout:"shared/productLayout.html",
		Name:"product/view.html",
		Data:iris.Map{
			"product" : product,
		},
	}
}

func (p *ProductController) GetOrder() []byte  {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productId, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userId, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 创建消息体
	message := datamodels.NewMessage(userId, productId)
	// 类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return []byte("true")

	//product, err := p.ProductService.GetProductById(int64(productId))
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}

	// 判断商品数量是否满足需求
	/*
	var orderId int64
	showMessage := "抢购失败"
	if product.ProductNum > 0 {
		// 扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		// 创建订单
		userId, err := strconv.Atoi(userId)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		
		order := &datamodels.Order{
			UserId:      int64(userId),
			ProductId:   int64(productId),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderId, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		showMessage = "抢购成功"
	}*/
}

