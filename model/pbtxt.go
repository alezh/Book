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
			bookCover := make([]library.BooksCache,0)
			dbmgo.PaginateNotSort("BookCover",bson.M{},i,pageSize,&bookCover)
			for _,p := range bookCover{
				value := library.ValueCache{Id:p.Id,Title:p.Title,Author:p.Author,CatalogUrl:p.CatalogUrl,Catalog:len(p.Catalog),Desc:p.Desc}
				pb.public.Add(p.Author+p.Title,-1,value)
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
			case "GetBook":
				go pb.getBook(k.Node)
			case "BookCover":
				go pb.bookCover(k.Node)
			case "getChapterDown":
				go pb.getChapterDown(k.Node)
			case "UpdateChapter":
				go pb.UpdateChapter(k.Node)
			default:
				fmt.Println("数据丢失",v)
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
			for _,v := range value.Data().(library.ValueCache).CatalogUrl {
				if v.Name == "pbtxt" {
					key = true
				}
			}
			if!key{
				orignalUrl := new(library.OriginUrl)
				orignalUrl.Name = "pbtxt"
				orignalUrl.Url  = pb.Web + url
				update := bson.M{"$push":bson.M{"catalogurl":orignalUrl}}
				dbmgo.UpdateSync("BookCover",value.Data().(library.ValueCache).Id,update)
			}
			pb.MQueue.InsertQueue(pb.Web + url,"UpdateChapter","")
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
//更新章节
func (pb *PbTxt)UpdateChapter(h *goquery.Document){
	Author,_ := h.Find("meta[name$=author]").Attr("content")
	Title ,_ := h.Find("meta[property$=title]").Attr("content")
	if value,err := pb.public.Value(Author + Title);err == nil{
		//获取章节
		aryClass := make([]interface{},0)
		Catalog  := make([]bson.ObjectId,0)
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
				chapter.Author = Author
				chapter.Title  = Title
				chapter.Group  = group
				chapter.ChapterName = selection.Text()
				chapter.Site = append(chapter.Site,Url)
				aryClass = append(aryClass,chapter)
				Catalog = append(Catalog,objId)
			}
		})
		if len(aryClass)>0{
			dbmgo.InsertAllSync("Chapter",aryClass...)
		}
		update := bson.M{"$push":bson.M{"Catalog":Catalog}}
		dbmgo.UpdateSync("BookCover",value.Data().(library.ValueCache).Id,update)
	}
}


//获取章节内容
func (pb *PbTxt)GetChapter(bookId string){
	Catalog := new(library.Catalog)
	Chapter := make([]library.SaveChapter,0)
	where := bson.M{"_id":bson.ObjectIdHex(bookId)}
	dbmgo.Finds("BookCover",where,Catalog)
	whereIn := bson.M{"_id": bson.M{"$in": Catalog.Catalog}}
	dbmgo.FindAllSort("Chapter",whereIn,"+sort",&Chapter)
	for _,v := range Chapter {
		pb.cache.Add(v.Site[0].Url,-1,v.Id)
		pb.MQueue.InsertQueue(v.Site[0].Url,"getChapterDown","")
	}
}
func (pb *PbTxt)getChapterDown(h *goquery.Document){
	if value,err :=pb.cache.Value(pb.WapWeb +h.Url.Path);err==nil{
		txt := h.Find("#nr1").Text()
		txtX := strings.TrimSpace(txt)
		content := new(library.UpChapterTxt)
		content.Content = strings.Replace(txtX, "\n\n    ", "\n", -1)
		update := bson.M{"$set":content}
		//fmt.Println(value.Data().(bson.ObjectId),update)
		dbmgo.UpdateSync("Chapter",value.Data().(bson.ObjectId),update)
	}


}