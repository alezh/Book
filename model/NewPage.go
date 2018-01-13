package model

import (
	"fmt"
	Tcp "Book/HttpConn"
	"golang.org/x/net/html"
	lib "Book/library"
)
//US站获取更新列表数据
type NewPage struct{
	Url string
	Title string
	NewChapter string
	Chapter string
}

type Node struct{
	Document *html.Node
}

type PageChan chan NewPage

var (
	Channel =make(chan *NewPage,500)
)

func init(){
	
}


func GetPage(url string){
	var Selection lib.Regexh
	doc := Tcp.GetNode(url)
	Selection.Selection =doc.Find("tr a")
	html :=Selection.ReAttr("title").AttrVal("href")
	fmt.Println(html)
	return
}