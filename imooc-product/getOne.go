package main

import "sync"

var sum int64 = 0

// 预存商品数量
var productNum int64 = 10000

// 互斥锁
var mutex sync.Mutex

// 获取秒杀商品
func GetOneProduct() bool {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()

	// 判断数据是否超限
	if sum < productNum {
		sum += 1
		return true
	}

	return false
}

