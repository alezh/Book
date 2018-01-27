package Thread

import (
	"time"
	"Book/HttpConn"
	"sync"
	"github.com/PuerkitoBio/goquery"
)

//控制网络链接数量

type MQueue struct {
	TotalThread   int
	Waiting       int
	NewThread       int
	Timer         time.Duration
	WaitingChan   chan int
	CounterChan   chan int
	ReduceChan    chan int
	OnlineTask    chan map[string]interface{}
	WrongChan     map[string]interface{}
	SuccessChan   chan map[string]interface{}
	Queue         map[string]interface{}
	WaitGroup     *sync.WaitGroup
}

type Wrong struct {
	Method string
	Count  int
	Assist  interface{}
}

type ThreadsNum struct {
	mu sync.Mutex
	TotalThread   int
	NewThread     int
}
type Response struct {
	Node    *goquery.Document
	Assist  interface{}
}
type Master struct {

}

func NewMQueue(num int , WaitGroup *sync.WaitGroup) *MQueue{
	rmq := new(MQueue)
	rmq.TotalThread = num      //链接总数
	rmq.WaitGroup = WaitGroup
	rmq.NewThread   = 0        //现有棕色
	rmq.Waiting     = 0
	rmq.WaitingChan = make(chan int)
	rmq.CounterChan = make(chan int)
	rmq.ReduceChan  = make(chan int)
	rmq.Timer       = time.Second * 10
	rmq.WrongChan   = make(map[string]interface{})
	rmq.SuccessChan = make(chan map[string]interface{})
	rmq.Queue       = make(map[string]interface{})
	go rmq.timerTask()
	return rmq
}

//TODO::插入列队
func (x *MQueue)InsertQueue(url string,method string,assist interface{}){
	//TODO::链接数 ++
	x.NewThread++

	if x.NewThread > x.TotalThread{
		//TODO::等待列队
		x.waiting()
	}

	//TODO::go数 ++
	x.WaitGroup.Add(1)
	go x.runTask(url,method,assist)
}

func (x *MQueue)waiting(){
	<- x.WaitingChan
}

//计数器
func (x *MQueue)counter()  {
	for{
		select {
		case <-x.CounterChan:
			x.NewThread++
		case <-x.ReduceChan:
			x.NewThread--
		}
	}
}
//执行任务
func (x *MQueue)runTask(url string,method string,assist interface{})  {
	if value,ok := HttpConn.HttpRequest(url);ok{
		if value != nil{
			var m = make(map[string]interface{})
			m[method]= &Response{value,assist}
			x.WaitGroup.Add(1)
			x.SuccessChan <- m
		}
	}else{
		if v,ok := x.WrongChan[url];ok{
			if sdk,o:= v.(*Wrong);o{
				if sdk.Count >= 40 {
					//抛弃错误
					delete(x.WrongChan,url)
				}else{
					x.WrongChan[url] = &Wrong{Count:sdk.Count+1,Method:method,Assist:assist}
				}
			}
		}else{
			x.WrongChan[url] = &Wrong{Count:1,Method:method}
		}
	}
	//TODO::连接数 --
	x.NewThread--

	//TODO:: 列队满时候 先进先出
	if x.NewThread >= x.TotalThread {
		x.WaitingChan <- 1
	}

	x.WaitGroup.Done()
}

//定时任务
func (x *MQueue)timerTask(){
	ticker := time.NewTicker(x.Timer)
	for{
		select {
		case <- ticker.C:
			x.wrongToQueue()
		}
	}
}

//报错任务重新插入列队
func (x *MQueue)wrongToQueue(){
	for k,v := range x.WrongChan{
		if sdk,o:= v.(*Wrong);o{
			delete(x.WrongChan,k)
			x.InsertQueue(k,sdk.Method,sdk.Assist)
		}
	}
}

