package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"imooc-shop/services"
)

type ProductController struct {
	Ctx				iris.Context
	ProductService 	services.ProductService
}

func (p *ProductController) GetAll() mvc.View  {
	product, err := p.ProductService.GetAllProduct()

	if err != nil {
		fmt.Println(err)
		fmt.Println()
	}

	return mvc.View{
		Name : "product/view.html",
		Data : iris.Map{
			"productArray" : product,
		},
	}
}

func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName:           "imooc",
	})

	if err := dec.Decode(p.Ctx.Request().Form, product);err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 跳转
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View  {
	return mvc.View{
		Name:"product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{
		TagName:           "imooc",
	})

	if err := dec.Decode(p.Ctx.Request().Form, product);err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	_, err := p.ProductService.InsertProduct(product)
	fmt.Println(err)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 跳转
	p.Ctx.Redirect("/product/all")
}



