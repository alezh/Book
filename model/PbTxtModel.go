package model

import (
	"Book/Thread"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"Book/library"
	"Book/dbmgo"
	"sync"
	"time"
	"strings"
	"fmt"
	"math"
	"gopkg.in/mgo.v2/bson"
)

type PbTxtModel struct {
	WebUrl            string
	LastUpUrl         string
	NewCreateUrl      string
	UnDesc            string
	MQueue            *Thread.MQueue
	NewBookPageSize   int
	WaitGroup         *sync.WaitGroup

}

func NewPbModel(wait *sync.WaitGroup)*PbTxtModel{
	pb := new(PbTxtModel)
	pb.WaitGroup    = wait
	pb.WebUrl       = "http://m.pbtxt.com"
	pb.LastUpUrl    = "http://m.pbtxt.com/top-lastupdate-"
	pb.NewCreateUrl = "http://m.pbtxt.com/top-postdate-"
	pb.MQueue       = Thread.NewMQueue(20,wait)
	pb.NewBookPageSize       = 0
	pb.UnDesc       = "最新章节推荐地址"
	go pb.receiving()
	return pb
}
//初始化 数据库 抓取书本
func (pb *PbTxtModel)Main(){
	pb.NewBook()
}

func (pb *PbTxtModel)getSqlToChapter(){
	count := dbmgo.Count("BookCover")
	pageSize := 4
	//向上取整
	key := int(math.Ceil(float64(count)/float64(pageSize)))
	for i:=1;i<=key ;i++ {
		var bookCover []library.BookCover
		dbmgo.Paginate("BookCover",bson.M{},"-created",i,pageSize,&bookCover)
		for _,p := range bookCover{
			pb.MQueue.InsertQueue(p.CatalogUrl.Url,"ChapterTxt",p)
		}
	}
}

//开始获取新书
func (pb *PbTxtModel)NewBook(){
	if pb.NewBookPageSize == 0{
		//获取页码
		pb.MQueue.InsertQueue(pb.NewCreateUrl + "1/","setCreatePageSize",nil)
		pb.NewBookPageSize = -1
	}else if pb.NewBookPageSize > 0{
		for a := 1;a<=pb.NewBookPageSize; a++{
			pb.MQueue.InsertQueue(pb.NewCreateUrl + strconv.Itoa(a) + "/","NewBook",nil)
		}
	}
}



//TODO::返回的数据接收数据
func (pb *PbTxtModel)receiving(){
	var f func(map[string]interface{})

	f = func(m map[string]interface{}) {
		for v,k := range m{
			value := k.(*Thread.Response)
			switch v {
			case "NewBook":
				pb.WaitGroup.Add(1)
				go pb.getNewBook(value.Node)
				pb.WaitGroup.Done()
				break
			case "BookCover":
				pb.WaitGroup.Add(1)
				go pb.getBookCover(value.Node)
				pb.WaitGroup.Done()
				break
			case "ChapterPage":
				pb.WaitGroup.Add(1)
				go pb.getChapterPage(value.Node)
				pb.WaitGroup.Done()
				break
			case "getChapterUrl":
				pb.WaitGroup.Add(1)
				go pb.getChapterUrl(value.Node)
				pb.WaitGroup.Done()
				break
			case "Chapter":
				pb.WaitGroup.Add(1)
				go pb.getChapter(value.Node,value.Assist)
				pb.WaitGroup.Done()
				break
			case "setCreatePageSize":
				pb.WaitGroup.Add(1)
				go pb.setCreatePageSize(value.Node)
				pb.WaitGroup.Done()
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
	pb.WaitGroup.Done()
	return
}
//获取书本
func (pb *PbTxtModel)getNewBook(doc *goquery.Document){
	doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
		class := new(library.Classify)
		class.Name = getStringNameZero("[",selection.Text(),"]")
		class.Author = getStringName("/",selection.Text(),"")
		html := selection.Find("a")
		class.Title = html.Text()
		url,_ := html.Attr("href")
		class.Url = pb.WebUrl + url
		dbmgo.InsertSync("Classify",class)
		//发送下载书本封面的请求
		pb.BookCover(class.Url)
	})
	pb.WaitGroup.Done()
	return
}
//封面
func (pb *PbTxtModel)BookCover(url string){
	pb.MQueue.InsertQueue(url,"BookCover",nil)
	return
}
//获取封面
func (pb *PbTxtModel)getBookCover(doc *goquery.Document)  {
	orignalUrl := new(library.OriginalUrl)
	bookCover := new(library.SaveBookCover)
	hTitle := doc.Find("title").Text()
	bookCover.Author = getString("(",hTitle,")_")
	bookCover.Title = getStringName("",hTitle,"(")
	bookCover.Status = "连载中"
	if coverImg ,err := doc.Find("div .block_img2 img").Attr("src");err{
		bookCover.CoverImg = coverImg
	}
	//bookCover.Sort = classify.Name
	//bookCover.Desc , _ = doc.Find("meta[name=description]").Attr("content")
	bookCover.Desc = getStringName("",doc.Find("div .intro_info").Text(),pb.UnDesc)
	orignalUrl.Name = "pbtxt"
	orignalUrl.Url = pb.WebUrl + doc.Url.Path + "page-1.html"
	bookCover.CatalogUrl = orignalUrl
	bookCover.Created = time.Now().Unix()
	dbmgo.InsertSync("BookCover",bookCover)
	pb.Chapter(orignalUrl.Url)
	pb.WaitGroup.Done()
	return
}
//章节
func (pb *PbTxtModel)Chapter(url string){
	pb.MQueue.InsertQueue(url,"ChapterPage",nil)
}
//章节分页
func (pb *PbTxtModel)getChapterPage(doc *goquery.Document){
	sel:=doc.Find(".listpage").First().Find("option")

	for i := range sel.Nodes{
		single := sel.Eq(i)
		if i > 0{
			if u ,e :=single.Attr("value");e{
				pb.MQueue.InsertQueue(pb.WebUrl + u,"getChapterUrl",nil)
			}
		}else{
			pb.getChapterUrl(doc)
		}
	}
	pb.WaitGroup.Done()
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
		dbmgo.InsertSync("ChapterUrl",&Url)
		//pb.MQueue.InsertQueue(Url.Url,"Chapter",Url)
		pb.WaitGroup.Done()
	})
	return
}
//下载章节
func (pb *PbTxtModel)getChapter(doc *goquery.Document,ass interface{}){
	chap := new(library.Chapter)
	if assist,ok := ass.(*library.OriginalUrl);ok{
		chap.Title = assist.Title
		chap.Url = assist.Url
		chap.Author = assist.Author
		chap.Sort = assist.Number
		chap.ChapterName = assist.Name
	}else {
		fmt.Println("assist nil",pb.WebUrl + doc.Url.Path )
	}
	chap.Content = pb.chapterTxt(doc.Find("#nr1").Text())
	chap.Site   = "pbtxt"
	dbmgo.InsertSync("Chapter",&chap)
	pb.WaitGroup.Done()
	return
}
//截取内容
func (pb *PbTxtModel)chapterTxt(txt string) (content string) {
	txtX := strings.TrimSpace(txt)
	content = strings.Replace(txtX, "\n\n    ", "\n", -1)
	return
}

func (pb *PbTxtModel)SqlToChapter(doc *goquery.Document,ass interface{}){

}