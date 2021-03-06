package main

import (
	"runtime"
	"Book/httprouter"
	"Book/controller"
	"Book/dbmgo"
	"sync"
	"net/http"
	"fmt"
	"Book/model"
)
const (
	WebUrl = "http://www.23us.so/top/lastupdate_"
)
func init() {	
    runtime.GOMAXPROCS(runtime.NumCPU()) // 多核多线程
	//mongoDb 初始化
	dbmgo.Init("127.0.0.1",27017,"BookDb")
	//dbmysql.Init()
}


var (
	waitGroup sync.WaitGroup
)

func main(){
	go Gather()
	router := Router()
	http.ListenAndServe(":8080", router)
}

func Router() *httprouter.Router{
	var bookRack controller.BookRack
	var index controller.Index
	var chapter controller.Chapter
	router := httprouter.New()
	router.OPTIONS("/Create/:os", index.Create)
	router.OPTIONS("/Save/:id/:book", bookRack.Save)
	router.OPTIONS("/Books/:id", bookRack.List)
	router.OPTIONS("/ChapterList/:bookId/:Site", chapter.ChapterList)
	router.OPTIONS("/Chapter/:chapterId", chapter.GetChapter)
	return router
}

func Gather(){
	//start := time.Now()
	//model.NewPbModel(&waitGroup).Main()
	model.NewPbTxt(&waitGroup).Main()
	//model.NewPbTxt(&waitGroup).Download()
	//model.NewXus(&waitGroup).Main()
	//model.NewXus(&waitGroup).UpdateBookCover()
	waitGroup.Wait()
	fmt.Print("采集完毕")
	//bp := new(PbTxt.BbLogic)
	//bp.Main()
	//bp.ChapterToNodes()
	//fmt.Println(bookCover[0].Id.Hex())
	//fmt.Printf("longCalculation took this amount of time: %s\n", time.Now().Sub(start))
}
