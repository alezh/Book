package main

import (
	"runtime"
	"Book/dbmgo"
)
const (
	WebUrl = "http://www.23us.so/top/lastupdate_"
)
func init() {	
    runtime.GOMAXPROCS(runtime.NumCPU()) // 多核多线程
	//mongoDb 初始化
	dbmgo.Init("127.0.0.1",27017,"BookDb")
}


func main(){
	//bp := new(PbTxt.BbLogic)
	//bp.Main()
	//bp.ChapterToNodes()


}