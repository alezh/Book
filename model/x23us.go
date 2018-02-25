package model

import (
	"gopkg.in/mgo.v2/bson"
	"Book/dbmgo"
	"Book/library"
	"Book/Thread"
	"Book/Cache"
	"sync"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"math"
)

type Xus struct {
	Web          string
	WapWeb       string
	WebClass     string
	WebBook      string
	MQueue      *Thread.MQueue
	WaitGroup   *sync.WaitGroup
	cache       *Cache.CacheTable
	public      *Cache.CacheTable
}
//有访问限制
func NewXus(wait *sync.WaitGroup)*Xus{
	us := new(Xus)
	us.Web          = "https://www.x23us.com"
	us.WapWeb       = "https://m.x23us.com"
	us.WebClass     = "https://m.x23us.com/class/"
	us.WebBook      = "https://m.x23us.com/book/"
	us.WaitGroup    = wait
	us.cache        = Cache.Create("x23us")
	us.public       = Cache.Create("public")
	us.init()
	us.MQueue       = Thread.NewMQueue(15,wait,"gbk")

	go us.receiving()
	return us
}
func (xu *Xus)receiving(){
	var f func(map[string]*Thread.Response)

	f = func(m map[string]*Thread.Response) {
		for v,k := range m{
			switch v {
			case "NewBook":
				go xu.getNewBook(k.Node)
				break
			case "getBook":
				go xu.classify(k.Node)
				break
			case "BookCover":
				go xu.bookCover(k.Node)
				break
			case "mBookCover":
				go xu.mBookCover(k.Node)
				break
			case "UpBookCover":
				go xu.upBookCover(k.Node)

			default:
				fmt.Println("数据丢失",v)
				break
			}
		}
	}
	for {
		value := <- xu.MQueue.SuccessChan
		f(value)
	}
}
func (xu *Xus)init(){
	//初始化
	fmt.Println("初始化数据缓存")
	if count := dbmgo.Count("BookCover");count>0{
		pageSize := 100
		//向上取整
		key := int(math.Ceil(float64(count)/float64(pageSize)))
		for i:=1;i<=key ;i++ {
			bookCover := make([]library.Books,0)
			dbmgo.PaginateNotSort("BookCover",bson.M{},i,pageSize,&bookCover)
			//fmt.Println("初始化:",int(math.Ceil(float64(pageSize*i)/float64(count)*100)),"%")
			for _,p := range bookCover{
				xu.public.Add(p.Author+p.Title,-1,p)
			}
		}
	}
}
//查询bookcover库
func (xu *Xus)Main(){
	xu.NewBook()
}

//10大分类
func (xu *Xus)NewBook(){
	fmt.Println("开始更新数据")
	xu.WaitGroup.Add(1)
	for i:=1;i<=10 ;i++  {
		url :=xu.WebClass + strconv.Itoa(i) + "_" + "1.html"
		xu.MQueue.InsertQueue(url,"NewBook","")
	}
}
//每个分类页数
func (xu *Xus)getNewBook(h *goquery.Document){
	if count,err := strconv.Atoi(getStringName("1/",h.Text(),"页)"));err ==nil {
		class := getStringName("/",h.Url.Path,"_")
		for i:=1;i<=count ;i++  {
			url := xu.WebClass + class + "_" + strconv.Itoa(i) +".html"
			xu.MQueue.InsertQueue(url,"getBook","")
		}
	}
}

func (xu *Xus)classify(h *goquery.Document){
	h.Find(".line").Each(func(_ int, selection *goquery.Selection){
		html := selection.Text()
		Title := getStringNameZero("",html,"/")
		Author := getStrings("/",html,"/")
		url ,_ := selection.Find("a").First().Attr("href")
		if value,err := xu.public.Value(Author + Title);err == nil{
			key := false
			for _,v := range value.Data().(library.Books).CatalogUrl {
				if v.Name == "x23us" {
					key = true
				}
			}
			if!key{
				orignalUrl := new(library.OriginUrl)
				orignalUrl.Name = "x23us"
				orignalUrl.Url  = xu.Web + url
				update := bson.M{"$push":bson.M{"catalogurl":orignalUrl}}
				dbmgo.UpdateSync("BookCover",value.Data().(library.Books).Id,update)
			}
		} else {
			//xu.MQueue.InsertQueue(xu.Web + url,"BookCover","")
			xu.MQueue.InsertQueue(xu.WapWeb + url,"mBookCover","")
		}
	})
}

func (xu *Xus)bookCover(h *goquery.Document){
	//获取封面
	orignalUrl := new(library.OriginUrl)
	bookCover := new(library.SaveBookCover)
	bookCover.Author,_ = h.Find("meta[name$=author]").Attr("content")
	bookCover.Title ,_ = h.Find("meta[property$=title]").Attr("content")
	bookCover.Status ,_ = h.Find("meta[name$=status]").Attr("content")
	bookCover.Sort  ,_ = h.Find("meta[name$=category]").Attr("content")
	if img,err:=h.Find("meta[property$=image]").Attr("content");err{
		if i := strings.LastIndex(img,"nocover");i<0{
			bookCover.CoverImg = img
		}
	}
	bookCover.NewChapter  ,_ = h.Find("meta[name$=latest_chapter_name]").Attr("content")
	desc ,_:= h.Find("meta[property$=description]").Attr("content")
	bookCover.Desc = strings.TrimSpace(desc)
	orignalUrl.Name = "x23us"
	orignalUrl.Url = xu.Web + h.Url.Path
	bookCover.CatalogUrl = append(bookCover.CatalogUrl,orignalUrl)
	//获取章节
	aryClass := make([]interface{},0)
	group := ""
	h.Find("#at tr").Each(func(i int, selection *goquery.Selection) {
		if k := selection.Find("th").Text();k!=""{
			group = k
		}
		selection.Find("a").Each(func(x int, s *goquery.Selection) {
			chapter := new(library.SaveChapter)
			Url := new(library.OriginUrl)
			objId := bson.NewObjectId()
			chapter.Id = objId
			chapter.Sort = i*4 + x +1
			Url.Name = "x23us"
			if cUrl,k :=s.Attr("href");k{
				Url.Url  = xu.WapWeb +h.Url.Path + cUrl
			}
			chapter.Author = bookCover.Author
			chapter.Title  = bookCover.Title
			chapter.Group  = group
			chapter.ChapterName = s.Text()
			chapter.Site = append(chapter.Site,Url)
			aryClass = append(aryClass,chapter)
			bookCover.Catalog = append(bookCover.Catalog,objId)
		})
	})
	if len(aryClass)>0{
		dbmgo.InsertAllSync("Chapter",aryClass...)
	}
	dbmgo.InsertSync("BookCover",bookCover)
}
//手机端没有限制
func (xu *Xus)mBookCover(h *goquery.Document){
	originlUrl := new(library.OriginUrl)
	bookCover := new(library.SaveBookCover)
	cover := h.Find(".index_block p")
	bookCover.Author = getStrings("：",cover.Eq(0).Text(),"")
	bookCover.Title = h.Find(".index_block h1").Text()
	bookCover.Status = "连载中"
	bookCover.Sort  = cover.Eq(1).Find("a").Text()
	originlUrl.Name = "x23us"
	originlUrl.Url = xu.WapWeb + h.Url.Path
	bookCover.CatalogUrl = append(bookCover.CatalogUrl,originlUrl)

	aryClass := make([]interface{},0)
	lens := len(h.Find(".chapter li").Nodes)

	h.Find(".chapter li").Each(func(i int, selection *goquery.Selection) {
		chapter := new(library.SaveChapter)
		Url := new(library.OriginUrl)
		objId := bson.NewObjectId()
		chapter.Id = objId
		chapter.Sort = lens - i
		Url.Name = "x23us"
		if cUrl,k :=selection.Find("a").Attr("href");k{
			Url.Url  = xu.WapWeb +h.Url.Path + cUrl
		}
		chapter.Author = bookCover.Author
		chapter.Title  = bookCover.Title
		chapter.ChapterName = selection.Text()
		chapter.Site = append(chapter.Site,Url)
		aryClass = append(aryClass,chapter)
		bookCover.Catalog = append(bookCover.Catalog,objId)
	})
	if len(aryClass)>0{
		dbmgo.InsertAllSync("Chapter",aryClass...)
	}
	dbmgo.InsertSync("BookCover",bookCover)
}

//更新封面
func (xu *Xus)UpdateBookCover(){
	xu.public.Foreach(func(key interface{}, item *Cache.CacheItem) {
		value := item.Data().(library.Books)
		for _,v := range value.CatalogUrl {
			if v.Name == "x23us" {
				id:= getStrings("/",getStrings("/html/",v.Url,""),"/")
				xu.MQueue.InsertQueue(xu.WebBook+ id,"UpBookCover","")
			}
		}
	})
}

//func (xu *Xus)upBookCover(h *goquery.Document){
//	//BookCover :=new(library.UpdateBookCover)
//	CoverImg:=""
//	if cUrl,k :=h.Find("block_img2 img").Attr("src");k{
//		if i := strings.LastIndex(cUrl,"nocover");i<0{
//			CoverImg = cUrl
//		}
//	}
//	html := h.Find(".block_txt2 p")
//
//	Title := h.Find(".block_txt2 h1").Text()
//	Author := html.Eq(2).Find("a").Text()
//	Sort := html.Eq(3).Find("a").Text()
//	Status:= getStrings("：",html.Eq(4).Text(),"")
//	NewChapter:=html.Eq(6).Find("a").Text()
//	desc := strings.TrimSpace(h.Find(".intro_info").Text())
//	update := bson.M{"$set":bson.M{"Sort":Sort,"Status":Status,"NewChapter":NewChapter}}
//	if "header" != desc{
//		update = bson.M{"$set":bson.M{"Sort":Sort,"Status":Status,"NewChapter":NewChapter,"Desc":desc}}
//		if CoverImg != ""{
//			update = bson.M{"$set":bson.M{"Sort":Sort,"Status":Status,"NewChapter":NewChapter,"Desc":desc,"CoverImg":CoverImg}}
//		}
//	}else if CoverImg != ""{
//		update = bson.M{"$set":bson.M{"Sort":Sort,"Status":Status,"NewChapter":NewChapter,"CoverImg":CoverImg}}
//	}
//	if value,err := xu.public.Value(Author + Title);err == nil{
//		fmt.Println(value,update)
//		//dbmgo.UpdateSync("BookCover",value.Data().(library.Books).Id,update)
//	}
//
//
//}

func (xu *Xus)upBookCover(h *goquery.Document){
	BookCover :=new(library.UpdateBookCover)
	//if cUrl,k :=h.Find("block_img2 img").Attr("src");k{
	//	if i := strings.LastIndex(cUrl,"nocover");i<0{
	//		BookCover.CoverImg = cUrl
	//	}
	//}
	html := h.Find(".block_txt2 p")

	Title := h.Find(".block_txt2 h1").Text()
	Author := html.Eq(2).Find("a").Text()
	BookCover.Sort = html.Eq(3).Find("a").Text()
	BookCover.Status= getStrings("：",html.Eq(4).Text(),"")
	BookCover.NewChapter=html.Eq(6).Find("a").Text()
	//desc := strings.TrimSpace(h.Find(".intro_info").Text())
	//if "header" != desc{
	//	BookCover.Desc = desc
	//}
	if value,err := xu.public.Value(Author + Title);err == nil{
		update := bson.M{"$set":BookCover}
		dbmgo.UpdateSync("BookCover",value.Data().(library.Books).Id,update)
	}


}
