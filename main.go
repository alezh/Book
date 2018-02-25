package main

import (
	"runtime"
	"Book/httprouter"
	"Book/controller"
	"Book/dbmgo"
	"sync"
	"net/http"
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
	router := Router()
	http.ListenAndServe(":8080", router)
	//start := time.Now()
	//model.NewPbModel(&waitGroup).Main()
	//model.NewPbTxt(&waitGroup).Main()
	//model.NewXus(&waitGroup).Main()
	//model.NewXus(&waitGroup).UpdateBookCover()
	//waitGroup.Wait()
	//bp := new(PbTxt.BbLogic)
	//bp.Main()
	//bp.ChapterToNodes()
	//fmt.Println(bookCover[0].Id.Hex())
	//fmt.Printf("longCalculation took this amount of time: %s\n", time.Now().Sub(start))

}

func Router() *httprouter.Router{
	var bookRack controller.BookRack
	var index controller.Index
	var chapter controller.Chapter
	router := httprouter.New()
	router.GET("/Create/:os", index.Create)
	router.GET("/Save/:id/:book", bookRack.Save)
	router.GET("/Books/:id", bookRack.List)
	router.GET("/ChapterList/:bookId/:Site", chapter.ChapterList)
	router.GET("/Chapter/:chapterId", chapter.GetChapter)
	return router
}
