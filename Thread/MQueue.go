package Thread

import (
	"time"
	"Book/HttpConn"
)

//控制网络链接数量

type MQueue struct {
	TotalThread   int
	Waiting       int
	Counter       int
	Timer         time.Duration
	WaitingChan   chan int
	CounterChan   chan int
	ReduceChan    chan int
	OnlineTask    chan map[string]interface{}
	WrongChan     map[string]interface{}
	SuccessChan   chan map[string]interface{}
	Queue         map[string]interface{}
}

type Wrong struct {
	Method string
	Count  int
}

type Master struct {

}

func NewMQueue(num int) *MQueue{
	rmq := new(MQueue)
	rmq.TotalThread = num
	rmq.Counter     = 0
	rmq.Waiting     = 0
	rmq.WaitingChan = make(chan int)
	rmq.CounterChan = make(chan int)
	rmq.ReduceChan  = make(chan int)
	rmq.Timer       = time.Second * 20
	rmq.WrongChan   = make(map[string]interface{})
	rmq.SuccessChan = make(chan map[string]interface{})
	rmq.Queue       = make(map[string]interface{})
	go rmq.counter()
	go rmq.timerTask()
	return rmq
}

func (x *MQueue)start(url string,method string){
	go x.runTask(url,method)
}
//TODO::插入列队
func (x *MQueue)InsertQueue(url string,method string){
	if x.Counter > x.TotalThread{
		//TODO::等待列队
		x.waiting()
	}
	//TODO::列队数 ++
	x.CounterChan <- 1
	go x.start(url,method)
}

func (x *MQueue)waiting(){
	<- x.WaitingChan
}

//计数器
func (x *MQueue)counter()  {
	for{
		select {
		case <-x.CounterChan:
			x.Counter++
		case <-x.ReduceChan:
			x.Counter--
		}
	}
}
//执行任务
func (x *MQueue)runTask(url string,method string)  {
	if value,ok := HttpConn.HttpRequest(url);ok{
		var m = make(map[string]interface{})
		m[method]=value
		x.SuccessChan <- m
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
	//TODO::列队数 --
	if x.Counter > x.TotalThread{
		x.ReduceChan <- 1
	}
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