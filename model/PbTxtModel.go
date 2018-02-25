package model

import (
	"Book/Thread"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"Book/library"
	"Book/dbmgo"
	"sync"
	"strings"
	"Book/Cache"
	"gopkg.in/mgo.v2/bson"
	"math"
	"fmt"
)

type PbTxtModel struct {
	Web               string
	WebUrl            string
	LastUpUrl         string
	NewCreateUrl      string
	UnDesc            string
	MQueue            *Thread.MQueue
	NewBookPageSize   int
	WaitGroup         *sync.WaitGroup
	cache             *Cache.CacheTable

}

func NewPbModel(wait *sync.WaitGroup)*PbTxtModel{
	pb := new(PbTxtModel)
	pb.WaitGroup    = wait
	pb.Web          = "http://www.pbtxt.com"
	pb.WebUrl       = "http://m.pbtxt.com"
	pb.LastUpUrl    = "http://m.pbtxt.com/top-lastupdate-"
	pb.NewCreateUrl = "http://m.pbtxt.com/top-postdate-"
	pb.MQueue       = Thread.NewMQueue(20,wait,"utf8")
	pb.NewBookPageSize       = 0
	pb.UnDesc       = "最新章节推荐地址"
	pb.cache        = Cache.Create("pbtxt")
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
		pb.MQueue.InsertQueue(pb.NewCreateUrl + "1/","setCreatePageSize","")
		pb.NewBookPageSize = -1
	}else if pb.NewBookPageSize > 0{
		for a := 1;a<=pb.NewBookPageSize; a++{
			pb.MQueue.InsertQueue(pb.NewCreateUrl + strconv.Itoa(a) + "/","NewBook","")
		}
	}
}

//TODO::返回的数据接收数据
func (pb *PbTxtModel)receiving(){
	var f func(map[string]*Thread.Response)

	f = func(m map[string]*Thread.Response) {
		for v,k := range m{
			switch v {
			case "NewBook":
				go pb.getNewBook(k.Node)
				break
			case "BookCover":
				go pb.getBookCover(k.Node)
				break
			case "ChapterPage":
				go pb.getChapterPage(k.Node)
				break
			case "getChapterUrl":
				go pb.getChapterUrl(k.Node)
				break
			case "Chapter":
				go pb.getChapter(k.Node)
				break
			case "setCreatePageSize":
				go pb.setCreatePageSize(k.Node)
				break
			case "ChapterPageDown":
				go pb.ChapterPageDown(k.Node)
				break
			case "ChapterUrlDown":
				go pb.ChapterUrlDown(k.Node)
				break
			case "ChapterDown":
				go pb.getChapterDown(k.Node)
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
	return
}
//获取书本
func (pb *PbTxtModel)getNewBook(doc *goquery.Document){
	//aryClass := make([]*library.MyClassify,0)
	aryClass := make([]interface{},0)
	doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
		//class := new(library.MyClassify)
		class := new(library.Classify)
		class.Name = getStringNameZero("[",selection.Text(),"]")
		class.Author = getStringName("/",selection.Text(),"")
		html := selection.Find("a")
		class.Title = html.Text()
		url,_ := html.Attr("href")
		class.Url = pb.WebUrl + url
		aryClass = append(aryClass,class)
		//发送下载书本封面的请求
		pb.BookCover(class.Url)
	})
	//dbmysql.Insert(aryClass)
	dbmgo.InsertAllSync("Classify",aryClass...)
	return
}
//封面
func (pb *PbTxtModel)BookCover(url string){
	pb.MQueue.InsertQueue(url,"BookCover","")
	return
}
//获取封面
func (pb *PbTxtModel)getBookCover(doc *goquery.Document)  {
	orignalUrl := new(library.OriginUrl)
	bookCover := new(library.SaveBookCover)
	//objId := bson.NewObjectId()
	//id := getStrings("/",doc.Url.Path,"/")
	//pb.cache.Add(id,-1,bson.NewObjectId())
	hTitle := doc.Find("title").Text()
	bookCover.Author = getString("(",hTitle,")_")
	bookCover.Title = getStringName("",hTitle,"(")
	bookCover.Status = "连载中"
	if coverImg ,err := doc.Find("div .block_img2 img").Attr("src");err{
		bookCover.CoverImg = coverImg
	}
	desc ,_:= doc.Find("meta[name=description]").Attr("content")
	bookCover.Desc = strings.TrimSpace(desc)
	//bookCover.Desc = getStringName("",doc.Find("div .intro_info").Text(),pb.UnDesc)
	orignalUrl.Name = "pbtxt"
	orignalUrl.Url = pb.WebUrl + doc.Url.Path + "page-1.html"
	bookCover.CatalogUrl = append(bookCover.CatalogUrl,orignalUrl)
	//bookCover.CatalogUrl = pb.WebUrl + doc.Url.Path + "page-1.html"
	//bookCover.Id = objId
	//dbmysql.Insert(bookCover)
	dbmgo.InsertSync("BookCover",bookCover)
	//pb.Chapter(orignalUrl.Url)
	return
}


//章节
func (pb *PbTxtModel)Chapter(url string){
	pb.MQueue.InsertQueue(url,"ChapterPage","")
}
//章节分页
func (pb *PbTxtModel)getChapterPage(doc *goquery.Document){
	sel:=doc.Find(".listpage").First().Find("option")

	for i := range sel.Nodes{
		single := sel.Eq(i)
		if i > 0{
			if u ,e :=single.Attr("value");e{
				pb.MQueue.InsertQueue(pb.WebUrl + u,"getChapterUrl","")
			}
		}else{
			pb.getChapterUrl(doc)
		}
	}
	return
}
//获取章节
func (pb *PbTxtModel)getChapterUrl(doc *goquery.Document){
	doc.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
		orUrl := getUrl(selection)
		Url := new(library.OriginalUrl)
		hTitle := doc.Find("title").Text()
		Url.Author = getString("(",hTitle,")_")
		Url.Title = getStringName("",hTitle,"(")
		Url.Name  = orUrl.Name
		Url.Number = orUrl.Number
		Url.Url = pb.WebUrl + orUrl.Url
		pb.cache.Add(orUrl.Url,-1,Url)
		pb.MQueue.InsertQueue(Url.Url,"Chapter","")
	})
	return
}
//下载章节
func (pb *PbTxtModel)getChapter(doc *goquery.Document){
	chap := new(library.Chapter)
	id := getStrings("/",doc.Url.Path,"/")
	//LOOK:
	if objId,err := pb.cache.Value(id);err == nil{
		chap.CoverId = objId.Data().(bson.ObjectId)
	}else{
		fmt.Println("id nil",id ,err.Error())
	}
	res, err := pb.cache.Value(doc.Url.Path)
	if err == nil{
		assist := res.Data().(*library.OriginalUrl)
		chap.Title = assist.Title
		chap.Url = assist.Url
		chap.Author = assist.Author
		chap.Sort = assist.Number
		chap.ChapterName = assist.Name
	}else {
		fmt.Println("assist nil",doc.Url.Path ,err.Error())
	}
	chap.Content = pb.chapterTxt(doc.Find("#nr1").Text())
	chap.Site   = "pbtxt"
	dbmgo.InsertSync("Chapter",&chap)
	return
}
//截取内容
func (pb *PbTxtModel)chapterTxt(txt string) (content string) {
	txtX := strings.TrimSpace(txt)
	content = strings.Replace(txtX, "\n\n    ", "\n", -1)
	return
}

//查询bookcover库
func (pb *PbTxtModel)SelectBookCover(){
	count := dbmgo.Count("BookCover")
	pageSize := 20
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key ;i++ {
		var bookCover []library.Books
		dbmgo.PaginateNotSort("BookCover",bson.M{},i,pageSize,&bookCover)
		//for _,p := range bookCover{
		//	id := getStrings("com/",p.CatalogUrl.Url,"/")
		//	pb.cache.Add(id,-1,p.Id)
		//	//pb.MQueue.InsertQueue(p.CatalogUrl.Url,"ChapterPageDown","")
		//}
	}
}
func (pb *PbTxtModel)ChapterPageDown(doc *goquery.Document){
	sel:=doc.Find(".listpage").First().Find("option")
	for i := range sel.Nodes{
		single := sel.Eq(i)
		if i > 0{
			if u ,e :=single.Attr("value");e{
				pb.MQueue.InsertQueue(pb.WebUrl + u,"ChapterUrlDown","")
			}
		}else{
			pb.ChapterUrlDown(doc)
		}
	}
	return
}
//获取章节
func (pb *PbTxtModel)ChapterUrlDown(doc *goquery.Document){
	doc.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
		orUrl := getUrl(selection)
		Url := new(library.OriginalUrl)
		hTitle := doc.Find("title").Text()
		Url.Author = getString("(",hTitle,")_")
		Url.Title = getStringName("",hTitle,"(")
		Url.Name  = orUrl.Name
		Url.Number = orUrl.Number
		Url.Url = pb.WebUrl + orUrl.Url
		pb.cache.Add(orUrl.Url,-1,Url)
		pb.MQueue.InsertQueue(Url.Url,"ChapterDown","")
	})
	return
}
//下载章节
func (pb *PbTxtModel)getChapterDown(doc *goquery.Document){
	chap := new(library.Chapter)
	id := getStrings("/",doc.Url.Path,"/")
	//LOOK:
	if objId,err := pb.cache.Value(id);err == nil{
		chap.CoverId = objId.Data().(bson.ObjectId)
	}else{
		fmt.Println("id nil",id ,err.Error())
	}
	res, err := pb.cache.Value(doc.Url.Path)
	if err == nil{
		assist := res.Data().(*library.OriginalUrl)
		chap.Title = assist.Title
		chap.Url = assist.Url
		chap.Author = assist.Author
		chap.Sort = assist.Number
		chap.ChapterName = assist.Name
	}else {
		fmt.Println("assist nil",doc.Url.Path ,err.Error())
	}
	chap.Content = pb.chapterTxt(doc.Find("#nr1").Text())
	chap.Site   = "pbtxt"
	dbmgo.InsertSync("Chapter",&chap)
	return
}

//----------------------------------------------------------------------------------------------------------------------
