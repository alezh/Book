package Thread

import (
	"time"
	"Book/HttpConn"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"fmt"
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
	SuccessChan   chan map[string]*goquery.Document
	Queue         map[string]interface{}
	WaitGroup     *sync.WaitGroup
}

type Wrong struct {
	Method string
	Count  int
}

type ThreadsNum struct {
	mu sync.Mutex
	TotalThread   int
	NewThread     int
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
	rmq.Timer       = time.Second * 20
	rmq.WrongChan   = make(map[string]interface{})
	rmq.SuccessChan = make(chan map[string]*goquery.Document)
	rmq.Queue       = make(map[string]interface{})
	go rmq.counter()
	go rmq.timerTask()
	return rmq
}

//func (x *MQueue)start(url string,method string){
//	x.WaitGroup.Add(1)
//	go x.runTask(url,method)
//}

//TODO::插入列队
func (x *MQueue)InsertQueue(url string,method string){
	//TODO::链接数 ++
	x.CounterChan <- 1
	if x.NewThread > x.TotalThread{
		//TODO::等待列队
		fmt.Println("现有连接数",x.NewThread)
		x.waiting()
	}
	//TODO::列队数 ++
	x.WaitGroup.Add(1)
	go x.runTask(url,method)
}

func (x *MQueue)waiting(){

	<- x.WaitingChan
	//for {
	//	if x.NewThread < x.TotalThread{
	//		break
	//	}
	//	select {
	//	case <- x.WaitingChan:
	//		break
	//	}
	//}
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
func (x *MQueue)runTask(url string,method string)  {
	if value,ok := HttpConn.HttpRequest(url);ok{
		if value != nil{
			var m = make(map[string]*goquery.Document)
			m[method]=value
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
					x.WrongChan[url] = &Wrong{Count:sdk.Count+1,Method:method}
				}
			}
		}
	}
	x.ReduceChan <- 1

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
			x.InsertQueue(k,sdk.Method)
		}
	}
}

