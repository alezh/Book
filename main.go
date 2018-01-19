package main

import (
	"runtime"
	"Book/dbmgo"
	"Book/model"
)
const (
	WebUrl = "http://www.23us.so/top/lastupdate_"
	Table = "BookInfo"
	chapter = "Chapter"
)
func init() {	
    runtime.GOMAXPROCS(runtime.NumCPU()) // 多核多线程
}

func init(){
	//mongoDb 初始化
	dbmgo.Init("127.0.0.1",27017,"BookDb")
}

func main(){
	//写入分类也书籍
	//info := model.PbTxtInfo()
	//info.GetSort()
	model.GetCover()
	//for a := 1;; a++{
	//	url := WebUrl + strconv.Itoa(a) + ".html"
	//	model.GetPage(url)
	//	break
	//}
}