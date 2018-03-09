package controller

import (
	"net/http"
	"Book/httprouter"
	"gopkg.in/mgo.v2/bson"
	"Book/dbmgo"
	"Book/library"
)

type Chapter struct {

}

//获取章节内容
func (c *Chapter)GetChapter(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	chapterId := params.ByName("chapterId")
	//siteName := params.ByName("Site")
	Content :=new(library.Content)
	dbmgo.Find("Chapter","_id",bson.ObjectIdHex(chapterId),Content)
	library.Render(w,Content,"","")
}
//获取章节列表
func (c *Chapter)ChapterList(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	bookId := params.ByName("bookId")
	siteName := params.ByName("Site")
	Catalog := new(library.Catalog)
	Chapter := make([]library.ChapterList,0)
	where := bson.M{"_id":bson.ObjectIdHex(bookId)}
	dbmgo.Finds("BookCover",where,Catalog)
	for k,v:=range Catalog.Site{
		if v.Name == siteName && k>0{
			//TODO::需要网站抓取章节
		}else {
			//TODO::自源
			whereIn := bson.M{"_id": bson.M{"$in": Catalog.Catalog}}
			dbmgo.FindAllSort("Chapter",whereIn,"+sort",&Chapter)
			library.Render(w,Chapter,"","")
			return
		}
	}
}

//获取换源
func (c *Chapter)GetSiteList(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
	bookId := params.ByName("bookId")
	Site :=new(library.Site)
	dbmgo.Find("BookCover","_id",bson.ObjectIdHex(bookId),Site)
	library.Render(w,Site,"","")
}
