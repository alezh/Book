package model

import (
	"Book/Thread"
	"sync"
	"Book/Cache"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"Book/dbmgo"
	"Book/library"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"math"
)

type PbTxt struct {
	Web          string
	WapWeb       string
	WebNBook       string
	MQueue      *Thread.MQueue
	WaitGroup   *sync.WaitGroup
	cache       *Cache.CacheTable
	public      *Cache.CacheTable
}

func NewPbTxt(wait *sync.WaitGroup)*PbTxt{
	pb :=new(PbTxt)
	pb.Web          = "http://www.pbtxt.com"
	pb.WapWeb       = "http://m.pbtxt.com"
	pb.WebNBook     = "http://m.pbtxt.com/top-postdate-"
	pb.cache        = Cache.Create("pbtxt")
	pb.public       = Cache.Create("public")
	pb.init()
	pb.MQueue       = Thread.NewMQueue(20,wait,"utf8")
	pb.WaitGroup    = wait
	go pb.receiving()
	return pb
}
func (pb *PbTxt)init(){
	//初始化
	fmt.Println("初始化数据缓存")
	if count := dbmgo.Count("BookCover");count>0{
		pageSize := 100
		//向上取整
		key := int(math.Ceil(float64(count)/float64(pageSize)))
		for i:=1;i<=key ;i++ {
			bookCover := make([]library.Books,0)
			dbmgo.PaginateNotSort("BookCover",bson.M{},i,pageSize,&bookCover)
			for _,p := range bookCover{
				pb.public.Add(p.Author+p.Title,-1,p)
			}
		}
	}
}
func (pb *PbTxt)receiving(){
	var f func(map[string]*Thread.Response)

	f = func(m map[string]*Thread.Response) {
		for v,k := range m{
			switch v {
			case "newBook":
				go pb.newBook(k.Node)
				break
			case "GetBook":
				go pb.getBook(k.Node)
				break
			case "BookCover":
				go pb.bookCover(k.Node)
				break
			default:
				fmt.Println("数据丢失",v)
				break
			}
		}
	}
	for {
		value := <- pb.MQueue.SuccessChan
		f(value)
	}
}
//查询bookcover库
func (pb *PbTxt)Main(){
	//pb.MQueue.InsertQueue("http://www.pbtxt.com/98380/","BookCover","")
	pb.MQueue.InsertQueue(pb.WebNBook+"1/","newBook","")
}
func (pb *PbTxt)newBook(h *goquery.Document){
	if count,err := strconv.Atoi(getStringName("1/",h.Text(),"页)"));err ==nil {
		pb.getBook(h)
		for i:=2;i<=count ;i++  {
			pb.MQueue.InsertQueue(pb.WebNBook + strconv.Itoa(i) + "/","GetBook","")
		}
	}
}
func (pb *PbTxt)getBook(h *goquery.Document){
	h.Find(".line").Each(func(_ int, selection *goquery.Selection){
		//Name := getStringNameZero("[",selection.Text(),"]")
		Author := getStringName("/",selection.Text(),"")
		html := selection.Find("a")
		Title := html.Text()
		url,_ := html.Attr("href")
		if value,err := pb.public.Value(Author + Title);err == nil{
			key := false
			for _,v := range value.Data().(library.Books).CatalogUrl {
				if v.Name == "pbtxt" {
					key = true
				}
			}
			if!key{
				orignalUrl := new(library.OriginUrl)
				orignalUrl.Name = "pbtxt"
				orignalUrl.Url  = pb.Web + url
				update := bson.M{"$push":bson.M{"catalogurl":orignalUrl}}
				dbmgo.UpdateSync("BookCover",value.Data().(library.Books).Id,update)
			}
		}else{
			pb.MQueue.InsertQueue(pb.Web + url,"BookCover","")
		}
	})
}
func (pb *PbTxt)bookCover(h *goquery.Document){
	//获取封面
	originUrl := new(library.OriginUrl)
	bookCover := new(library.SaveBookCover)
	bookCover.Author,_ = h.Find("meta[name$=author]").Attr("content")
	bookCover.Title ,_ = h.Find("meta[property$=title]").Attr("content")
	bookCover.Status = "连载中"
	bookCover.Sort  ,_ = h.Find("meta[name$=category]").Attr("content")
	if img,err:=h.Find("meta[property$=image]").Attr("content");err{
		if i := strings.LastIndex(img,"nocover");i<0{
			bookCover.CoverImg = img
		}
	}
	bookCover.Desc = h.Find(".intro p").Eq(1).Text()
	originUrl.Name = "pbtxt"
	originUrl.Url = pb.Web + h.Url.Path
	bookCover.CatalogUrl = append(bookCover.CatalogUrl,originUrl)
	//获取章节
	aryClass := make([]interface{},0)
	group := ""
	h.Find("dt,dd").Each(func(i int, selection *goquery.Selection){
		if x := strings.LastIndex(selection.Text(),"》");x>=0{
			group = selection.Text()
		}else{
			chapter := new(library.SaveChapter)
			Url := new(library.OriginUrl)
			objId := bson.NewObjectId()
			chapter.Id = objId
			chapter.Sort = i+1
			Url.Name = "pbtxt"
			if cUrl,k :=selection.Find("a").Attr("href");k{
				Url.Url  = pb.WapWeb +h.Url.Path + cUrl
			}
			chapter.Author = bookCover.Author
			chapter.Title  = bookCover.Title
			chapter.Group  = group
			chapter.ChapterName = selection.Text()
			chapter.Site = append(chapter.Site,Url)
			aryClass = append(aryClass,chapter)
			bookCover.Catalog = append(bookCover.Catalog,objId)
		}
	})
	if len(aryClass)>0{
		dbmgo.InsertAllSync("Chapter",aryClass...)
	}
	dbmgo.InsertSync("BookCover",bookCover)
}