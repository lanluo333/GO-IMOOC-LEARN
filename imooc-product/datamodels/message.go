package datamodels

// 简单的消息体
type Message struct {
	ProductId	int64
	UserId		int64
}

// 创建结构体
func NewMessage(userId, productId int64) *Message {
	return &Message{
		ProductId: productId,
		UserId:    userId,
	}
}
