package main

import (
	"runtime"
	"Book/dbmgo"
	"Book/model"
	"sync"
)
const (
	WebUrl = "http://www.23us.so/top/lastupdate_"
)
func init() {	
    runtime.GOMAXPROCS(runtime.NumCPU()) // 多核多线程
	//mongoDb 初始化
	dbmgo.Init("127.0.0.1",27017,"BookDb")
}


var (
	waitGroup sync.WaitGroup
)

func main(){
	//var index controller.Index
	//router := httprouter.New()
	//router.GET("/", index.Index)
	//http.ListenAndServe(":8080", router)
	//start := time.Now()
	model.NewPbModel(&waitGroup).Main()
	waitGroup.Wait()
	//bp := new(PbTxt.BbLogic)
	//bp.Main()
	//bp.ChapterToNodes()
	//fmt.Println(bookCover[0].Id.Hex())
	//fmt.Printf("longCalculation took this amount of time: %s\n", time.Now().Sub(start))

}