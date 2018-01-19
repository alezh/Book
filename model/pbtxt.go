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
	chanBookCover = make(chan *library.BookCover,10)
	coverOne = make(chan []*library.Sort,0)
	coverTwo = make(chan []*library.Sort,0)
	coverThree = make(chan []*library.Sort,0)
	coverFour = make(chan []*library.Sort,0)
	thread  = 4
	UnDesc = "最新章节推荐地址"

}

type Pbtxt struct{
	WebUrl string
	Sort   []*Classification
	CacheSort chan *library.Sort
}

type Classification struct{
	Class string
	Name  string
}

func PbTxtInfo() (info Pbtxt){
	info.WebUrl = webUrl
	info.Sort = Sort
	info.CacheSort = make(chan *library.Sort)
	go info.sortSave()
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
	var err error
	for _,v := range info.Sort{
		sum := 1
		for a := 0;a<=sum; a++{
			doc := HttpConn.GetSelect(webUrl + v.Class + strconv.Itoa(a) + ".html")
			var Selection library.Regexh
			if sum == 1{
				if sum , err = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); err != nil{
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
					info.CacheSort <- &sort
				})
		}
	}
	close(info.CacheSort)
	defer Distinct()
	return
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

func (info *Pbtxt)sortSave(){
	for {
		v := <- info.CacheSort
		dbmgo.InsertToDB("Sort",v)
	}
}

func bookCoverSave(){
	for {
		v := <- chanBookCover
		dbmgo.InsertToDB("BookCover",v)
	}
}

//多线程叠加抓取 大流量 N+N
//func bookCoverSave(){
//	for {
//		select{
//		case v:=<- chanBookCover :
//			dbmgo.InsertToDB("BookCover",v)
//		case a := <- coverOne:
//			go logicCover(a)
//		case b := <- coverTwo:
//			go logicCover(b)
//		case c := <- coverThree:
//			go logicCover(c)
//		case d := <- coverFour:
//			go logicCover(d)
//		}
//	}
//}

func bookCoverOne(){
	for {
		v := <- coverOne
		logicCover(v)
	}
}
func bookCoverTwo(){
	for {
		v := <- coverTwo
		logicCover(v)
	}
}
func bookCoverThree(){
	for {
		v := <- coverThree
		logicCover(v)
	}
}
func bookCoverFour(){
	for {
		v := <- coverFour
		logicCover(v)
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
		if len(sort) == 0 {
			break
		}
	}
}

//分页去除分类书籍
func GetCover(){
	fmt.Println("抓取书籍封面")
	go bookCoverSave()
	go bookCoverOne()
	go bookCoverTwo()
	go bookCoverThree()
	go bookCoverFour()
	count := dbmgo.Count("SortOnly")
	pageSize := 100
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key;i++{
		var Sort []*library.Sort
		dbmgo.Paginate("SortOnly",bson.M{},"-count",i,pageSize,&Sort)
		multipleThread(Sort)
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
		var bookCover library.BookCover
		var orignalUrl library.OriginalUrl
		doc := HttpConn.GetSelect(value.Url)
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
	}
}

func logicChapter()  {

}




