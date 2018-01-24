package Thread

type RequestManage struct {
	CoreOne   chan []interface{}
	CoreTwo   chan []interface{}
	CoreThree chan []interface{}
	CoreFour  chan []interface{}
	CoreFive  chan []interface{}
	CoreSix   chan []interface{}
	CoreSeven chan []interface{}
	CoreEight chan []interface{}
	CacheSort chan []interface{}
	CoreNina  chan []interface{}
	WaitingOne chan int
	WaitingTwo chan int
	WaitingThree chan int
	WaitingFour chan int
	WaitingFive chan int
	WaitingSix chan int
	WaitingSeven chan int
	WaitingEight chan int
}



type Master struct {

}