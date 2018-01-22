package PbTxt

import (
	"Book/library"
	"fmt"
	"Book/dbmgo"
	"Book/HttpConn"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"math"
)


type BbLogic struct{
	Class []*library.Classify
	Thread int                            //线程数量
	UnDesc string                         //简介剔除字段
	CoreOne chan []interface{}
	CoreTwo chan []interface{}
	CoreThree chan []interface{}
	CoreFour chan []interface{}
	CoreFive chan []*library.Sort
	CoreSix chan []*library.Sort
	CoreSeven chan []*library.Sort
	CoreEight chan []*library.Sort
	CacheSort chan *library.Sort
	CoreNina chan *library.Sort
	SaveDb chan *library.SaveDb           //数据保存
	Test chan []interface{}
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
	}
}

//多线程叠加抓取 大流量 N+N
func (v *BbLogic)BookCoverSave(){
	for {
		select{
		case o := <- v.SaveDb:
			dbmgo.InsertAllSync(o.Table,o.Data)
			 v.CountChannel <- 1
		case a := <- v.CoreOne:
			go downloadCover(a)
		case b := <- v.CoreTwo:
			go logicCover(b)
		case c := <- v.CoreThree:
			go logicCover(c)
		case d := <- v.CoreFour:
			go logicCover(d)
		case e := <- v.CoreFive:
			go logicCover(e)
		case f := <- v.CoreSix:
			go logicCover(f)
		case g := <- v.CoreSeven:
			go logicCover(g)
		case h := <- v.CoreEight:
			go logicCover(h)
		case i := <- v.CoreNina:
			go logicCover(i)
		case j := <- v.CacheSort:
			go logicCover(j)

		}
	}
}

func (v *BbLogic)Main(){
	v.CoreOne    = make(chan []interface{},100)
	v.CoreTwo    = make(chan []interface{},70)
	v.CoreThree  = make(chan []interface{},50)
	v.CoreFour   = make(chan []interface{},1)
	v.CoreFive   = make(chan []*library.Sort,1)
	v.CoreSix    = make(chan []*library.Sort,1)
	v.CoreSeven  = make(chan []*library.Sort,1)
	v.CoreEight  = make(chan []*library.Sort,1)
	v.CoreNina   = make(chan *library.Sort)
	v.CacheSort  = make(chan *library.Sort)
	v.SaveDb     = make(chan *library.SaveDb,1)
	v.Test       = make(chan []interface{})
	v.CountChannel = make(chan int )
	v.CountProgress     = 0
	v.TotalProgress     = 0
	go v.BookCoverSave()
	go v.cumulation()
	return
}

//获取新书
func (v *BbLogic)Classify(){
	sum := 1
	var erro error
	for a := 0;a<=sum; a++{
	LOOK:
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
			v.SaveDb <- Db  //20写入
			go v.Cover(aryClass)
			v.TotalProgress ++
			fmt.Println(v.TotalProgress)
		}else{
			fmt.Println("Goto Try again ")
			goto LOOK
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

//拆分数组给 channel
func (v *BbLogic)multipleThread(Classifyes []interface{}){
	length := len(Classifyes)
	key := int(math.Ceil(float64(length)/float64(v.Thread)))
	for a := 0;a<v.Thread; a++{
		start := a * key
		end := start + key
		if end>length{
			end = length
		}
		if start >length{
			fmt.Println("start > length ：" + strconv.Itoa(start))
			break
		}
		sort :=Classifyes[start:end]
		if len(sort) == 0 {
			break
		}
		switch a {
		case 0:
			v.CoreOne <- sort
			break
		case 1:
			v.CoreTwo <- sort
			break
		case 2:
			v.CoreThree <- sort
			break
		case 3:
			v.CoreFour <- sort
			break
		default:
			break
		}

	}
}


//分析书籍封面
func (v *BbLogic)Cover(sort []interface{}){
	for _,value := range sort {
		if s,ok := value.(*library.Classify);ok{
			fmt.Println(s.Url)
		}
	}
}
//下载书籍
func downloadCover(sort []interface{})  {
	for _,value := range sort {
		if s,ok := value.(*library.Classify);ok{
			fmt.Println(s.Url)
		}
	}
}