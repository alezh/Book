package Thread

import (
	"time"
	"Book/HttpConn"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"Book/Cache"
	"runtime"
)

//控制网络链接数量

type MQueue struct {
	TotalThread   int
	Waiting       int
	NewThread     int
	encoded       string
	Timer         time.Duration
	WaitingChan   chan int
	CounterChan   chan int
	ReduceChan    chan int
	OnlineTask    *Cache.CacheTable
	WrongChan     map[string]*Wrong
	SuccessChan   chan map[string]*Response
	Queue         map[string]*Response
	WaitGroup     *sync.WaitGroup
}

type Wrong struct {
	Method string
	Count  int
	Key    string
}

type Response struct {
	Node    *goquery.Document
	Key      string
}

func NewMQueue(num int , WaitGroup *sync.WaitGroup, encoded string) *MQueue{
	rmq := new(MQueue)
	rmq.TotalThread = num      //链接总数
	rmq.WaitGroup = WaitGroup
	rmq.NewThread   = 0        //现有棕色
	rmq.Waiting     = 0
	rmq.encoded     = encoded
	rmq.WaitingChan = make(chan int)
	rmq.CounterChan = make(chan int)
	rmq.ReduceChan  = make(chan int)
	rmq.Timer       = time.Second * 10
	rmq.OnlineTask  = Cache.Create("OnlineTask")
	rmq.WrongChan   = make(map[string]*Wrong)
	rmq.SuccessChan = make(chan map[string]*Response)
	rmq.Queue       = make(map[string]*Response)
	go rmq.timerTask()
	go rmq.timer()
	rmq.WaitGroup.Add(3)
	return rmq
}

//TODO::插入列队
func (x *MQueue)InsertQueue(url string,method string,cKey string){

	//TODO::链接数 ++
	x.NewThread++

	if x.NewThread > x.TotalThread{
		go x.counter(url,method,true)
		//TODO::等待列队
		x.waiting()
		go x.counter(url,method,false)
	}
	go x.runTask(url,method,cKey)
}


func (x *MQueue)waiting(){
	<- x.WaitingChan
}

//计数器
func (x *MQueue)counter(url string,method string,k bool)  {
	if k{
		x.OnlineTask.Add(url,0,method)
	}else{
		x.OnlineTask.Delete(url)
	}
	return
}
//执行任务
func (x *MQueue)runTask(url string,method string,cKey string)  {
	if value,ok := HttpConn.HttpRequest(url,x.encoded);ok{
		if value != nil{
			var m = make(map[string]*Response)
			m[method] = &Response{value,cKey}
			x.SuccessChan <- m
		}
	}else{
		if v,ok := x.WrongChan[url];ok{
			if v.Count >= 40 {
				//抛弃错误
				delete(x.WrongChan,url)
			}else{
				x.WrongChan[url] = &Wrong{Count:v.Count+1,Method:method,Key:cKey}
			}
		}else{
			x.WrongChan[url] = &Wrong{Count:1,Method:method}
		}
	}
	//TODO::连接数 --
	x.NewThread--

	//TODO::列队满时候 先进先出
	if x.NewThread >= x.TotalThread {
		x.WaitingChan <- 1
	}

	return
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
func (x *MQueue)timer()  {
	timer := time.NewTicker(time.Second * 5)
	for{
		select {
		case <- timer.C:
			x.log()
		}
	}
}

//报错任务重新插入列队
func (x *MQueue)wrongToQueue(){
	for k,v := range x.WrongChan{
		delete(x.WrongChan,k)
		go x.InsertQueue(k,v.Method,v.Key)
	}
}

func (x *MQueue)log(){
	green   := string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	reset   := string([]byte{27, 91, 48, 109})
	blue    := string([]byte{27, 91, 57, 55, 59, 52, 52, 109})

	fmt.Printf("Waiting task => %s %d %s | connections => %s %d %s | NumGoroutine => %s %d %s \n",green,x.OnlineTask.Count(),reset,blue,x.NewThread-x.OnlineTask.Count(),reset,blue,runtime.NumGoroutine(),reset)
	//捞出阻塞的数据
	if x.OnlineTask.Count()>0 && x.NewThread <= 0 {
		x.WaitingChan <- 1
	}else if x.OnlineTask.Count() == 0 && x.NewThread <= 0 {
		//循环3次 等待15秒
		x.WaitGroup.Done()
	}
}

