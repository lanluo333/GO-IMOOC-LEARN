package services

import (
	"imooc-shop/datamodels"
	"imooc-shop/repositories"
)

type IOrderService interface {
	GetOrderByID(int64)(*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo()(map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (orderId int64, err error)
}

func NewOrderService(OrderManagerRepository repositories.OrderManagerRepository) *OrderService  {
	return &OrderService{OrderRepository:OrderManagerRepository}
}

type OrderService struct {
	OrderRepository  repositories.OrderManagerRepository
}

func (o *OrderService) GetOrderByID(orderId int64)(*datamodels.Order, error) {
	return o.OrderRepository.SelectByKey(orderId)
}

func (o *OrderService) DeleteOrderByID(orderId int64) bool  {
	return o.OrderRepository.Delete(orderId)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error  {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error)  {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error)  {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error)  {
	return o.OrderRepository.SelectAllWithInfo()
}

// 根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderId int64, err error)  {
	order := &datamodels.Order{
		UserId:      message.UserId,
		ProductId:   message.ProductId,
		OrderStatus: datamodels.OrderSuccess,
	}

	return o.InsertOrder(order)
}

