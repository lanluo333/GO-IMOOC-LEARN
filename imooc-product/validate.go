package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"imooc-shop/encrypt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"imooc-shop/rabbitmq"
	"time"
)

// 设置集群地址，最好内网IP
var hostArray = []string{"192.168.1.6", "192.168.1.6"}

var localHost = ""

// 数量控制接口服务器内网IP，或者SLB或者内网IP
var GetOneIp  = "127.0.0.1"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.Consistent

// rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

// 用来存放控制信息
type AccessControl struct{
	// 用来存放用户想要存放的信息
	sourcesArray map[int]time.Time
	sync.RWMutex
}

// 服务器间隔时间
var interval = 20

// 创建全局变量
var accessControl = &AccessControl{sourcesArray:make(map[int]time.Time)}

// 获取指定数据
func (m *AccessControl) GetNewRecord(uid int) time.Time  {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool  {
	// 获取用户uid
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	// 采用一致性hash算法，根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	// 判断是否为本机
	if hostRequest == localHost {
		// 执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	}else {
		// 不是本机，充当代理去访问数据返回结果
		return m.GetDataFromOtherMap(hostRequest, req)
	}
}

// 黑名单
type BlackList struct {
	listArray	map[int]bool
	sync.RWMutex
}

// 获取黑名单
var blackList = &BlackList{
	listArray: make(map[int]bool),
}

// 添加黑名单
func (m *BlackList) SetBlackListByID(uid int) bool  {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}

func (m *BlackList) GetBlackListByID(uid int) bool  {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

// 处理本机map，并且处理业务逻辑，返回的结果类型为bool
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool)  {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	// 判断黑名单
	if blackList.GetBlackListByID(uidInt) {
		return false
	}

	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		// 业务判断，是否在指定的范围之后
		if dataRecord.Add(time.Duration(interval)*time.Second).After(time.Now()) {
			return false
		}
	}

	m.SetNewRecord(uidInt)
	return true
}

// 获取其他节点获取结果
func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool  {

	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}

	// 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		}else {
			return false
		}
	}

	return false
}

// 设置记录
func (m *AccessControl) SetNewRecord(uid int)  {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourcesArray[uid] = time.Now()
}

func CheckRight(w http.ResponseWriter, r *http.Request)  {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// 执行正常业务逻辑
func check(w http.ResponseWriter, r *http.Request)  {
	// 执行正常业务逻辑
	fmt.Println("执行check")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <=0 {
		w.Write([]byte("false,productID 获取失败"))
		return
	}

	productString := queryForm["productID"][0]
	fmt.Println(productString)

	// 获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false,用户 uid 获取失败"))
		return
	}

	// 1. 分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false,分布式权限验证出错"))
		return
	}

	// 2. 获取数量控制权限
	hostUrl := "http://"+GetOneIp+":"+GetOnePort+"/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false，抢单失败：" + err.Error()))
		return
	}
	// 判断控制数量接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// 整合下单
			productId, err := strconv.ParseInt(productString,10,64)
			if err != nil {
				w.Write([]byte("false"))
			}
			userId, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
			}
			// 3. 创建消息体
			message := datamodels.Message{ProductId:productId, UserId:userId}
			// 类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
			}

			// 4. 生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
			}
			w.Write([]byte("true"))
			return
		}
	}

	w.Write([]byte("false"))
	return
}

// 统一注册拦截器,每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error  {
	// 添加基于cookie的权限验证
	fmt.Println("cookie验证")
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	fmt.Println("执行验证")
	return nil
}

func CheckUserInfo(r *http.Request) error  {
	// 获取uid cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		fmt.Println(err)
		return errors.New("获取uid失败")
	}

	// 获取签名
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("获取用户签名失败")
	}

	// 对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return err
	}

	fmt.Println("用户uid：" + uidCookie.Value)
	fmt.Println("签名uid：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}

	return errors.New("身份校验失败")
}

// 自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool  {
	if checkStr == signStr {
		return true
	}

	return false
}

// 模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error)  {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}

	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET",hostUrl, nil)
	if err != nil {
		return
	}

	// 手动指定，排除多余cookie
	cookieUid := &http.Cookie{Name:"uid", Value:uidPre.Value}
	cookieSign := &http.Cookie{Name:"sign", Value:uidSign.Value}
	// 添加cookie到模拟请求
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// 获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(response.Body)

	return
}

func main()  {
	// 负载均衡器设置
	// 采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	// 采用一致性hash算法添加节点
	for _,v := range hostArray {
		hashConsistent.Add(v)
	}

	// 获取本机ip
	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localIp)

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()

	// 设置静态文件目录
	// 当访问localhost:xxxx/html，会路由到fileserver进行处理
	//当访问URL为/html/example.txt时，fileserver会将/html与URL进行拼接，得到/tmp/html/example.txt，
	// 而实际上example.txt的地址是/tmp/example.txt，因此这样将访问不到相应的文件，返回404 NOT FOUND。
	// 因此解决方案就是把URL中的/html/去掉，而http.StripPrefix做的就是这个。
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))
	// 设置资源目录
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./fronted/web/public"))))


	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilterUrl("/check", Auth)
	filter.RegisterFilterUrl("/checkRight", Auth)
	// 2.启动服务
	http.HandleFunc("/check",filter.Handle(check))
	http.HandleFunc("/checkRight",filter.Handle(CheckRight))
	http.ListenAndServe(":8083", nil)
}


