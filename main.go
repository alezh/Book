package main

import (
	"Book/model"
	"runtime"
	//"strconv"
	//"time"
	//"fmt"
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

}

func main(){
	model.PbTxtInfo()
	//for a := 1;; a++{
	//	url := WebUrl + strconv.Itoa(a) + ".html"
	//	model.GetPage(url)
	//	break
	//}
}