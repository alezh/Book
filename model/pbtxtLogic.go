package model

import (
	"Book/library"
	"fmt"
	"Book/dbmgo"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
	"Book/HttpConn"
)


type PbTxtLogic struct{
	Sort []*Classification
	ChanBookCover chan *library.BookCover //书本封面
	Thread int
	UnDesc string
	CoreOne chan []*library.Sort
	CoreTwo chan []*library.Sort
	CoreThree chan []*library.Sort
	CoreFour chan []*library.Sort
	CoreFive chan []*library.Sort
	CoreSix chan []*library.Sort
	CoreSeven chan []*library.Sort
	CoreEight chan []*library.Sort
	CacheSort chan *library.Sort
	CoreNina chan *library.Sort
	Test chan string
	CountChannel chan int
	CountX int
}

func init(){

}

//同步累计器
func (v *PbTxtLogic)cumulation()  {
	for{
		k:= <- v.CountChannel
		v.CountX = v.CountX + k
		fmt.Println(v.CountX)
	}
}

//多线程叠加抓取 大流量 N+N
func (v *PbTxtLogic)BookCoverSave(){
	for {
		select{
		case x:=<- v.ChanBookCover :
			dbmgo.InsertSync("BookCover",x)
			 v.CountChannel <- 1
		case w := <- v.CacheSort:
			dbmgo.InsertToDB("Sort",w)
			fmt.Println("insert")
			v.CountChannel <- 1
		case a := <- v.CoreOne:
			go logicCover(a)
		case b := <- v.CoreTwo:
			go logicCover(b)
		case c := <- v.CoreThree:
			go logicCover(c)
		case d := <- v.CoreFour:
			go logicCover(d)
		}
	}
}

func (v *PbTxtLogic)Main(){
	v.ChanBookCover = make(chan *library.BookCover,500)
	v.CoreOne    = make(chan []*library.Sort,1)
	v.CoreTwo    = make(chan []*library.Sort,1)
	v.CoreThree  = make(chan []*library.Sort,1)
	v.CoreFour   = make(chan []*library.Sort,1)
	v.CoreFive   = make(chan []*library.Sort,1)
	v.CoreSix    = make(chan []*library.Sort,1)
	v.CoreSeven  = make(chan []*library.Sort,1)
	v.CoreEight  = make(chan []*library.Sort,1)
	v.CoreNina   = make(chan *library.Sort)
	v.CacheSort  = make(chan *library.Sort)
	v.Test  = make(chan string)
	v.CountChannel = make(chan int )
	v.CountX = 0
	go v.BookCoverSave()
	return
}


func (v *PbTxtLogic)GetLastUpdate(){
	go v.cumulation()
	sum := 1
	var erro error
	for a := 0;a<=sum; a++{
		if doc,err := HttpConn.GetSelect(lastUpdate + strconv.Itoa(a) + "/");err{
			if sum == 1{
				if sum , erro = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); erro != nil{
					fmt.Println("页码获取错误")
					//暂停循环
					break
				}
			}
			doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
				var sort library.Sort
				sort.Name = getStringName("[",selection.Text(),"]")
				sort.Author = getStringName("/",selection.Text(),"")
				html := selection.Find("a")
				sort.Title = html.Text()
				url,_ := html.Attr("href")
				sort.Url = webUrl + url
				if !dbmgo.Finds("Sort",bson.M{"title":sort.Title,"author":sort.Author},&sort){
					v.CacheSort <- &sort
					fmt.Println("title"+sort.Title+"author"+sort.Author)
				}
			})
		}
	}
	for   {
		//查看通道任务是否完成
		if v.CountX == sum{
			break
		}
	}
	return
}

