package PbTxt

import (
	"Book/library"
	"fmt"
	"Book/dbmgo"
	"Book/HttpConn"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)


type BbLogic struct{
	Class []*library.Classify
	Thread int                            //线程数量
	UnDesc string                         //简介剔除字段
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
	SaveDb chan *library.SaveDb           //数据保存
	Test chan string
	CountChannel chan int                 //传输完成进度
	CountProgress int                     //完成进度
	TotalProgress int                     //总进度
}

func init(){
}

//同步累计器
func (v *BbLogic)cumulation()  {
	for{
		k:= <- v.CountChannel
		v.CountProgress = v.CountProgress + k
		fmt.Println(v.CountProgress)
	}
}

//多线程叠加抓取 大流量 N+N
func (v *BbLogic)BookCoverSave(){
	for {
		select{
		case o := <- v.SaveDb:
			dbmgo.InsertAllSync(o.Table,o.Data)
			 v.CountChannel <- 1
		//case a := <- v.CoreOne:
		//	go logicCover(a)
		//case b := <- v.CoreTwo:
		//	go logicCover(b)
		//case c := <- v.CoreThree:
		//	go logicCover(c)
		//case d := <- v.CoreFour:
		//	go logicCover(d)
		}
	}
}

func (v *BbLogic)Main(){
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
	v.SaveDb     = make(chan *library.SaveDb,1)
	v.Test       = make(chan string)
	v.CountChannel = make(chan int )
	v.CountProgress     = 0
	v.TotalProgress     = 0
	go v.BookCoverSave()
	go v.cumulation()
	return
}

func (v *BbLogic)Classify(){
	sum := 1
	var erro error
	for a := 0;a<=sum; a++{
		if doc,err := HttpConn.HttpRequest(newCreate + strconv.Itoa(a) + "/");err{
			if sum == 1{
				if sum , erro = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); erro != nil{
					fmt.Println("页码获取错误")
					break
				}
			}
			Db := new(library.SaveDb)
			Db.Table = "Classify"
			var aryClass []interface{}
			doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
				var class library.Classify
				class.Name = getStringNameZero("[",selection.Text(),"]")
				class.Author = getStringName("/",selection.Text(),"")
				html := selection.Find("a")
				class.Title = html.Text()
				url,_ := html.Attr("href")
				class.Url = webUrl + url
				aryClass = append(aryClass,&class)
			})
			Db.Data = aryClass
			v.SaveDb <- Db
			v.TotalProgress ++
			fmt.Println(v.TotalProgress)
		}else{
			break
		}
	}
	for{
		//查看通道任务是否完成
		if v.CountProgress == v.TotalProgress{
			break
		}
	}
	fmt.Println("close")
	return
}

