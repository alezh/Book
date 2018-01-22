package model

import (
	"Book/HttpConn"
	"Book/library"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"Book/dbmgo"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

const (
	//爬虫抓取 网站地址
	webUrl   = "http://m.pbtxt.com"
	lastUpdate ="http://m.pbtxt.com/top-lastupdate-"
	//分类
	// ["xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"]
	//玄幻|奇幻|修真|武侠|历史|都市|网游|科幻|恐怖|穿越|言情|校园
)

var Sort []*Classification
var chanBookCover chan *library.BookCover //书本封面
var thread int
var UnDesc string
var coverOne chan []*library.Sort
var coverTwo chan []*library.Sort
var coverThree chan []*library.Sort
var coverFour chan []*library.Sort
var CacheSort chan *library.Sort
var countChannel chan int
var countX int
var SaveDb chan *library.SaveDb

//初始化分类
func init(){
	name := []string{"xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"}
	sort := []string{"/xuanhuan/","/xiuzhen/","/wuxia/","/lishi/","/dushi/","/game/","/kehuan/","/kongbu/","/chuanyue/","/yanqing/","/xiaoyuan/"}
	for k,_ := range sort{
		var Class Classification
		Class.Name = name[k]
		Class.Class = sort[k]
		Sort = append(Sort,&Class)
	}
	chanBookCover = make(chan *library.BookCover,500)
	coverOne = make(chan []*library.Sort,0)
	coverTwo = make(chan []*library.Sort,0)
	coverThree = make(chan []*library.Sort,0)
	coverFour = make(chan []*library.Sort,0)
	CacheSort = make(chan *library.Sort)
	SaveDb = make(chan *library.SaveDb)
	thread  = 4
	UnDesc = "最新章节推荐地址"
	countChannel = make(chan int ,1)
	countX = 0

}

type Pbtxt struct{
	WebUrl string
	Sort   []*Classification
}

type Classification struct{
	Class string
	Name  string
}

func PbTxtInfo() (info Pbtxt){
	info.WebUrl = webUrl
	info.Sort = Sort
	go sortSave()
	return
}
//sort去重
func Distinct(){
	fmt.Println("书籍分类去重")
	dbmgo.Aggregation("Sort","SortOnly")
}
//logic——————————————————————————————————————————————————————————————————————————————————————————————————————————————

//获取分类
func (info *Pbtxt)GetSort(){
	fmt.Println("抓取书籍分类")
	var erro error
	for _,v := range info.Sort{
		sum := 1
		for a := 0;a<=sum; a++{
			if doc,err := HttpConn.HttpRequest(webUrl + v.Class + strconv.Itoa(a) + ".html");err{
				var Selection library.Regexh
				if sum == 1{
					if sum , erro = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); erro != nil{
						fmt.Println("页码获取错误")
						//暂停循环
						break
					}
				}
				Selection.Selection = doc.Find("p")
				Selection.Selection.Each(
					func(_ int, selection *goquery.Selection) {
						var sort library.Sort
						sort.Name = v.Name
						sort.Author = getStringName("/",selection.Text(),"")
						html := selection.Find("a")
						sort.Title = html.Text()
						url,_ := html.Attr("href")
						sort.Url = webUrl + url
						CacheSort <- &sort
					})
			}

		}
	}
	close(CacheSort)
	defer Distinct()
	return
}

func GetLastUpdate(){
	sum := 1
	var erro error
	for a := 0;a<=sum; a++{
		if doc,err := HttpConn.HttpRequest(lastUpdate + strconv.Itoa(a) + "/");err{
			if sum == 1{
				if sum , erro = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); erro != nil{
					fmt.Println("页码获取错误")
					//暂停循环
					break
				}
				doc.Find("div .cover p").Each(func(_ int, selection *goquery.Selection) {
					var sort library.Sort
					sort.Name = getStringName("[",selection.Text(),"]")
					sort.Author = getStringName("/",selection.Text(),"")
					html := selection.Find("a")
					sort.Title = html.Text()
					url,_ := html.Attr("href")
					sort.Url = webUrl + url
					var DbSort library.Sort
					if !dbmgo.Finds("Sort",bson.M{"title":sort.Title,"author":sort.Author},&DbSort){
						CacheSort <- &sort
					}
				})
			}
		}
	}
}

//截取字符串
func getStringName(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if f!= ""{
		if i := strings.LastIndex(text,f);i>0{
			if l != ""{
				if k := strings.LastIndex(text,l);k>0{
					c := text[i+len(f):k]
					return c
				}
			}
			c := text[i+len(f):]
			return c
		}
	}else{
		if k := strings.Index(text,l);k>0{
			c := text[:k]
			return c
		}
	}
	return text
}

func sortSave(){
	for {
		v := <- CacheSort
		dbmgo.InsertToDB("Sort",v)
	}
}

//多线程叠加抓取 大流量 N+N
func bookCoverSave(){
	for {
		select{
		case v:=<- chanBookCover :
			dbmgo.InsertSync("BookCover",v)
			countChannel <- 1
		case a := <- coverOne:
			go logicCover(a)
		case b := <- coverTwo:
			go logicCover(b)
		case c := <- coverThree:
			go logicCover(c)
		case d := <- coverFour:
			go logicCover(d)
		case o:= <-SaveDb:
			dbmgo.InsertSync(o.Table,o.Data)
			countChannel <- 1
		}
	}
}

//拆分数组给 channel
func multipleThread(Sort []*library.Sort){
	length := len(Sort)
	key := int(math.Ceil(float64(length)/float64(thread)))
	for a := 0;a<thread; a++{
		start := a * key
		end := start + key
		if end>length{
			end = length
		}
		if start >length{
			fmt.Println("start > length ：" + strconv.Itoa(start))
			break
		}
		sort :=Sort[start:end]
		if len(sort) == 0 {
			break
		}
		switch a {
		case 0:
			coverOne <- sort
			break
		case 1:
			coverTwo <- sort
			break
		case 2:
			coverThree <- sort
			break
		case 3:
			coverFour <- sort
			break
		default:
			break
		}

	}
}
//同步累计器
func cumulation()  {
	for{
		v:= <-countChannel
		countX = countX + v
	}
}

//分页去除分类书籍
func GetCover(){
	fmt.Println("抓取书籍封面")
	go bookCoverSave()
	go cumulation()
	count := dbmgo.Count("SortOnly")
	pageSize := 100
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key;i++{
		var Sort []*library.Sort
		dbmgo.Paginate("SortOnly",bson.M{},"-count",i,pageSize,&Sort)
		multipleThread(Sort)
	}
	for {
		//查看通道任务是否完成
		if countX == count{
			break
		}
		fmt.Println(countX)
	}
	defer close(coverOne)
	defer close(coverTwo)
	defer close(coverThree)
	defer close(coverFour)
	return
}

//下载书籍封面
func logicCover(Sort []*library.Sort){
	for _,value := range Sort{
		if doc,err := HttpConn.HttpRequest(value.Url);err{
			var bookCover library.BookCover
			var orignalUrl library.OriginalUrl
			bookCover.Author = getStringName("",value.Author,"&")
			bookCover.Title = value.Title
			bookCover.Status = "连载中"
			if coverImg ,err := doc.Find("div .block_img2 img").Attr("src");err{
				bookCover.CoverImg = coverImg
			}
			bookCover.Sort = value.Name
			bookCover.Desc = getStringName("",doc.Find("div .intro_info").Text(),UnDesc)
			orignalUrl.Name = value.Name
			orignalUrl.Url = value.Url + "page-1.html"
			bookCover.CatalogUrl = &orignalUrl
			bookCover.Created = time.Now().Unix()
			chanBookCover <- &bookCover
		}else{
			var reset []*library.Sort
			reset = append(reset,value)
			coverFour <- reset
			fmt.Println("->Try again channel")
		}
	}
}

//获取章节
func GetChapter()  {
	go bookCoverSave()
	go cumulation()
	count := dbmgo.Count("BookCover")
	pageSize := 100
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key;i++{
		var bookCover []*library.BookCover
		dbmgo.Paginate("BookCover",bson.M{},"-count",i,pageSize,&bookCover)
		logicChapter(bookCover)
	}
	for {
		//查看通道任务是否完成
		if countX == count{
			break
		}
		fmt.Println(countX)
	}
}

//获取章节
func logicChapter(cover []*library.BookCover){
	for _,value := range cover{
		var chap library.Chapter
		Catalogs := ChapterToNodes(value.CatalogUrl.Url)
		chap.Title = value.Title
		chap.CoverId = value.Id
		chap.Chapters= Catalogs
		var Db library.SaveDb
		Db.Table = "Catalog"
		Db.Data = chap
		SaveDb <- &Db
	}
}
//每本的章节
func ChapterToNodes(Url string) (Catalogs []*library.Catalog){
	if doc,err := HttpConn.HttpRequest(Url);err{
		sel:=doc.Find(".listpage").First().Find("option")
		for i := range sel.Nodes {
			single := sel.Eq(i)
			if i > 0{
				if u ,e :=single.Attr("value");e{
					if docs,errs := HttpConn.HttpRequest(webUrl +u);errs{
						docs.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
							html:=selection.Find("a")
							if v,b :=html.Attr("href");b{
								var cataLog library.Catalog
								cataLog.Url   = webUrl + v
								cataLog.Title = getStringName("、",html.Text(),"")
								cataLog.Content = ChapterTxt(cataLog.Url)
								Catalogs = append(Catalogs,&cataLog)
							}
						})
					}
				}
			}else{
				doc.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
					html:=selection.Find("a")
					if v,b :=html.Attr("href");b{
						var catlog library.Catalog
						catlog.Url   = webUrl + v
						catlog.Title = getStringName("、",html.Text(),"")
						catlog.Content = ChapterTxt(catlog.Url)
						Catalogs = append(Catalogs,&catlog)
					}
				})
			}
		}
	}
	return
}

func ChapterTxt(Url string)string{
	if doc,err := HttpConn.HttpRequest(Url);err{
		txt := doc.Find("#nr1").Text()
		return txt
	}
	return ""
}



