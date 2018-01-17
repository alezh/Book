package main

import (
	"runtime"
	"Book/dbmgo"
	"Book/library"
	"fmt"
	"gopkg.in/mgo.v2/bson"
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
	info := model.PbTxtInfo()
	info.GetSort()
	//读取数据测试
	//data := library.Sort{}
	//err :=dbmgo.Find("Sort","title","巫师进化手札",&data)
	//if err {
	//	fmt.Println(data)
	//}
	//for a := 1;; a++{
	//	url := WebUrl + strconv.Itoa(a) + ".html"
	//	model.GetPage(url)
	//	break
	//}
}