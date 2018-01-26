package model

import (
	"Book/Thread"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"Book/library"
	"Book/dbmgo"
	"sync"
)

type PbTxtModel struct {
	WebUrl         string
	LastUpUrl      string
	NewCreateUrl            string
	MQueue                  *Thread.MQueue
	NewBookPageSize         int
	WaitGroup      *sync.WaitGroup

}

func NewPbModel(wait *sync.WaitGroup)*PbTxtModel{
	pb := new(PbTxtModel)
	pb.WaitGroup    = wait
	pb.WebUrl       = "http://m.pbtxt.com"
	pb.LastUpUrl    = "http://m.pbtxt.com/top-lastupdate-"
	pb.NewCreateUrl = "http://m.pbtxt.com/top-postdate-"
	pb.MQueue       = Thread.NewMQueue(10,wait)
	pb.NewBookPageSize       = 0
	go pb.receiving()
	return pb
}
//初始化 数据库 抓取书本
func (pb *PbTxtModel)Main(){
	pb.NewBook()
}

//开始获取新书
func (pb *PbTxtModel)NewBook(){
	if pb.NewBookPageSize == 0{
		//获取页码
		pb.MQueue.InsertQueue(pb.NewCreateUrl + "1/","setCreatePageSize")
		pb.NewBookPageSize = -1
	}else if pb.NewBookPageSize > 0{
		for a := 1;a<=pb.NewBookPageSize; a++{
			pb.MQueue.InsertQueue(pb.NewCreateUrl + strconv.Itoa(a) + "/","NewBook")
		}
	}
}



//TODO::返回的数据接收数据
func (pb *PbTxtModel)receiving(){

	var f func(map[string]*goquery.Document)

	f = func(m map[string]*goquery.Document) {
		for v,k := range m{
			switch v {
			case "setCreatePageSize":
				pb.WaitGroup.Add(1)
				go pb.setCreatePageSize(k)
				pb.WaitGroup.Done()
			case "NewBook":
				pb.WaitGroup.Add(1)
				go pb.getNewBook(k)
				pb.WaitGroup.Done()
			case "BookCover":


			}
		}
	}
	for {
		 value := <- pb.MQueue.SuccessChan
		 f(value)
	}
}

//设置新增书本总翻页数
func (pb *PbTxtModel)setCreatePageSize(h *goquery.Document)  {
	if sum , err := strconv.Atoi(getStringName("1/",h.Text(),"页)")); err == nil{
		if sum > 0 {
			pb.NewBookPageSize = sum
			//设置完成后去下载书本
			pb.NewBook()
		}else{
			pb.NewBookPageSize = -1
		}
	}
	pb.WaitGroup.Done()
}

func (pb *PbTxtModel)getNewBook(doc *goquery.Document){
	doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
		var class library.Classify
		class.Name = getStringNameZero("[",selection.Text(),"]")
		class.Author = getStringName("/",selection.Text(),"")
		html := selection.Find("a")
		class.Title = html.Text()
		url,_ := html.Attr("href")
		class.Url = pb.WebUrl + url
		dbmgo.InsertSync("Classify",&class)
		pb.BookCover(&class)
	})
	pb.WaitGroup.Done()
}

func (pb *PbTxtModel)BookCover(book *library.Classify){
	pb.MQueue.InsertQueue(book.Url,"BookCover")
}

func (pb *PbTxtModel)getBookCover(doc *goquery.Document)  {
	
}