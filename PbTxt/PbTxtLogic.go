package PbTxt

import (
	"Book/library"
	"fmt"
	"Book/dbmgo"
	"Book/HttpConn"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"math"
	"time"
	"gopkg.in/mgo.v2/bson"
)


type BbLogic struct{
	Class []*library.Classify
	Thread int                            //线程数量
	UnDesc string                         //简介剔除字段
	CoreOne   chan []interface{}
	CoreTwo   chan []interface{}
	CoreThree chan []interface{}
	CoreFour  chan []interface{}
	CoreFive  chan []interface{}
	CoreSix   chan []interface{}
	CoreSeven chan []interface{}
	CoreEight chan []interface{}
	CacheSort chan []interface{}
	CoreNina  chan []interface{}
	SaveDb chan *library.SaveDb           //数据保存
	SaveCD chan *library.SaveDb           //保存缓存数据
	CacheDb map[string][]interface{}      //缓存数据
	CountChannel chan int                 //传输完成进度
	CountProgress int                     //完成进度
	TotalProgress int                     //总进度
	CacheSize int                         //sql缓存条数
	OnlineTask map[string]int             //在线运行的模块线程
	WaitingSeven chan int                           //等待协程
	WaitingEight chan int                           //等待协程
	WaitingTimeOut chan int                           //等待协程

}
func (v *BbLogic)Main(){
	v.CoreOne    = make(chan []interface{})
	v.CoreTwo    = make(chan []interface{})
	v.CoreThree  = make(chan []interface{})
	v.CoreFour   = make(chan []interface{})
	v.CoreFive   = make(chan []interface{})
	v.CoreSix    = make(chan []interface{})
	v.CoreSeven  = make(chan []interface{})
	v.CoreEight  = make(chan []interface{})
	v.CoreNina   = make(chan []interface{})
	v.CacheSort  = make(chan []interface{})
	v.SaveDb     = make(chan *library.SaveDb)
	v.SaveCD     = make(chan *library.SaveDb)
	v.CacheDb    = make(map[string][]interface{})
	v.CountChannel = make(chan int )
	v.CountProgress     = 0
	v.TotalProgress     = 0
	v.Thread     = 4
	v.UnDesc     = "最新章节推荐地址"
	v.CacheSize  = 60
	v.OnlineTask  = make(map[string]int)
	v.WaitingSeven  = make(chan int)
	v.WaitingEight  = make(chan int)
	v.WaitingTimeOut  = make(chan int)
	go v.BookCoverSave()
	go v.cumulation()
	return
}

//同步累计器
func (v *BbLogic)cumulation()  {
	for{
		select {
		    case <- v.CountChannel:
			    v.CountProgress ++
		}
	}
}

//多线程 N+N
func (v *BbLogic)BookCoverSave(){
	ticker := time.NewTicker( 20 * time.Second)
	for {
		select{
		case o := <- v.SaveDb:
			dbmgo.InsertAllSync(o.Table,o.Data...)
			v.CountChannel <- 1
		case a := <- v.CoreOne:
			go v.logicProcess(a)
		case b := <- v.CoreTwo:
			go v.logicProcess(b)
		case c := <- v.CoreThree:
			go v.logicProcess(c)
		case d := <- v.CoreFour:
			go v.logicProcess(d)
		case e := <- v.CoreFive:
			go v.logicProcess(e)
		case f := <- v.CoreSix:
			go v.logicProcess(f)
		//case g := <- v.CoreSeven:
		//	go v.logicProcess(g)
		//case h := <- v.CoreEight:
		//	go v.logicProcess(h)
		case i := <- v.CoreNina:
			go v.logicProcess(i)
		//case j := <- v.CacheSort:
		//	go v.logicProcess(j)
		 case <- ticker.C:
			go v.timerToDb()

		}
	}
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
			doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
				var class library.Classify
				class.Name = getStringNameZero("[",selection.Text(),"]")
				class.Author = getStringName("/",selection.Text(),"")
				html := selection.Find("a")
				class.Title = html.Text()
				url,_ := html.Attr("href")
				class.Url = webUrl + url
				Db.Data = append(Db.Data,&class)
			})
			v.TotalProgress ++
			v.SaveDb <- Db
			v.multipleThread(Db.Data)
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
	fmt.Println("Classify Close")
	return
}

//发送任务 倍数拆分数组给 channel
func (v *BbLogic)multipleThread(process []interface{}){
	length := len(process)
	v.TotalProgress += length //统计任务数量
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
		sort :=process[start:end]
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
		case 4:
			v.CoreFive <- sort
			break
		case 5:
			v.CoreSix <- sort
			break
		case 6:
			v.CoreSeven <- sort
			break
		case 7:
			v.CoreEight <- sort
			break
		default:
			break
		}

	}
}

//协程分配
func (v *BbLogic)doubleThread(process []interface{}){
	length := len(process)
	v.TotalProgress += length //统计任务数量
	key := int(math.Ceil(float64(length)/float64(2)))
	for a := 0;a<2; a++{
		start := a * key
		end := start + key
		if end>length{
			end = length
		}
		if start >length{
			fmt.Println("start > length ：" + strconv.Itoa(start))
			break
		}
		sort :=process[start:end]
		if len(sort) == 0 {
			break
		}
		switch a {
		case 0:
			v.singleTask(sort)
			v.WaitingSeven<-1
			break
		case 1:
			v.singleTask(sort)
			v.WaitingEight<-1
			break
		default:
			break
		}
	}
}

//逻辑处理分配任务
func (v *BbLogic)logicProcess(process ...interface{}){
	var Db library.SaveDb
	Save := true
	v.TotalProgress ++ //新增消费
	for _,value := range process {
		if sdk,ok := value.(*library.Classify);ok{
			//书本封面
			Db.Table = "BookCover"
			if value,k :=v.downloadCover(sdk);k{
				Db.Data = append(Db.Data,value)
			}
		}else if b,c := value.(*library.BookCover);c{
			Save = false
			var Dbs library.SaveDb
			Dbs.Table = "Chapter"
			value := v.downloadChapter(b)
			if len(value)>0 {
				Dbs.Data = value
				v.SaveDb <- &Dbs
			}
		}
	}
	if Save{
		v.SaveDb <- &Db
	}else{
		v.CountChannel <- 1
	}
}


func (v *BbLogic)singleTask(task []interface{}){
	for _,value := range task{
		if b,c := value.(*library.BookCover);c{
			var Dbs library.SaveDb
			Dbs.Table = "Chapter"
			value := v.downloadChapter(b)
			if len(value)>0 {
				Dbs.Data = value
				v.SaveDb <- &Dbs
			}
		}
	}
}

//下载书籍封面
func (v *BbLogic)downloadCover(classify *library.Classify) (bookCover *library.BookCover , ok bool) {
	if doc,err := HttpConn.HttpRequest(classify.Url);err{
		var orignalUrl library.OriginalUrl
		bookCover.Author = getStringName("",classify.Author,"&")
		bookCover.Title = classify.Title
		bookCover.Status = "连载中"
		if coverImg ,err := doc.Find("div .block_img2 img").Attr("src");err{
			bookCover.CoverImg = coverImg
		}
		bookCover.Sort = classify.Name
		bookCover.Desc = getStringName("",doc.Find("div .intro_info").Text(),v.UnDesc)
		orignalUrl.Name = classify.Name
		orignalUrl.Url = classify.Url + "page-1.html"
		bookCover.CatalogUrl = &orignalUrl
		bookCover.Created = time.Now().Unix()
		return bookCover,true
	}else{
		var reset []interface{}
		reset = append(reset,classify)
		v.CoreNina <- reset
		fmt.Println("->Try again channel")
		return nil,false
	}
}

//sql写入缓存
func (y *BbLogic)cacheToDb(db *library.SaveDb){
	if v,ok := y.CacheDb[db.Table];ok{
		if  l:=len(v);l >= y.CacheSize {
			//删除map缓存
			delete(y.CacheDb,db.Table)
			v = append(v,db.Data...)
			dbmgo.InsertAllSync(db.Table,v...)
			y.CountChannel <- l
		}else{
			//加入缓存数据
			y.CacheDb[db.Table] = append(v,db.Data...)
		}
	}else{
		//新增map缓存
		y.CacheDb[db.Table] = db.Data
	}
}

//定时清理没有写入的缓存sql
func (y *BbLogic)timerToDb(){
	//for k,v := range y.CacheDb{
	//	if _,ok := y.OnlineTask[k];ok{
	//		delete(y.CacheDb,k)
	//		dbmgo.InsertAllSync(k,v...)
	//		fmt.Println("Db Timer to mgo")
	//	}
	//}
	fmt.Println(y.TotalProgress)
}

//获取章节内容
func (y *BbLogic)ChapterToNodes(){
	count := dbmgo.Count("BookCover")
	pageSize := 2
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key ;i++  {
		var bookCover []*library.BookCover
		var inset []interface{}
		dbmgo.Paginate("BookCover",bson.M{},"-created",i,pageSize,&bookCover)
		for _,p := range bookCover{
			inset = append(inset,p)
		}
		go y.doubleThread(inset)
		<-y.WaitingSeven
		<-y.WaitingEight
		fmt.Println("next")
	}
}

func(v *BbLogic)downloadChapter(Book *library.BookCover)(Chapter []interface{}){
	var originalUrl []*library.OriginalUrl
	LOOK:
	if doc,err := HttpConn.HttpRequest(Book.CatalogUrl.Url);err{
		originalUrl =getChapter(doc)
	}else{
		goto LOOK
		return
	}
	for _, n := range originalUrl {
		var chap library.Chapter
		chap.Title = Book.Title
		//chap.CoverId = Book.Id
		chap.Url = n.Url
		chap.Author = Book.Author
		chap.ChapterName = n.Name
		chap.Content= chapterTxt(n.Url)
		chap.Sort = n.Number
		Chapter= append(Chapter,&chap)
	}
	return
}

