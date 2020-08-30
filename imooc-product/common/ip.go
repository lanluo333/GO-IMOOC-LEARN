package common

import (
	"errors"
	"net"
)

func GetIntranceIp() (string, error)  {
	addres, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addres {
		// 检查IP地址，判断是否回环地址
		// 將 address 做型別轉換成 *net.IPNet
		// 环bai回地址是主机用于向du自身发送通信的一个特zhi殊地址。
		// 环回dao地址为同一台设备上运行的 TCP/IP 应用程序和服务之间相互通信提供了一条捷径
		if ipnet, ok := address.(*net.IPNet);ok&&!ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("获取地址异常")
}


