package controller

import (
	"net/http"
	"Book/httprouter"
	"gopkg.in/mgo.v2/bson"
	"Book/dbmgo"
	c "Book/library"
)

//书架
type BookRack struct {
}

//书架列表
func (b *BookRack)List(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	userId := params.ByName("id")
	user   := new(c.UserDb)
	books  := make([]c.Book,0)
	dbmgo.Find("User","_id",bson.ObjectIdHex(userId),user)
	whereIn := bson.M{"_id": bson.M{"$in": user.BookList}}
	dbmgo.FindAll("BookCover",whereIn,&books)
	c.Render(w,books,"","")
}
//新增书本
func (b *BookRack)Save(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
	bookId := params.ByName("book")
	userId := params.ByName("id")
	update := bson.M{"$push":bson.M{"booklist":bson.ObjectIdHex(bookId)}}
	dbmgo.UpdateSync("User",bson.ObjectIdHex(userId),update)
	c.Render(w,"ok","","")
}
//删除书本
func (b *BookRack)Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	bookId := params.ByName("book")
	userId := params.ByName("id")
	update := bson.M{"$pop":bson.M{"booklist":bson.ObjectIdHex(bookId)}}
	dbmgo.UpdateSync("User",bson.ObjectIdHex(userId),update)
	c.Render(w,"ok","","")
}