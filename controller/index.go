package controller

import (
	"net/http"
	"Book/httprouter"
	c "Book/library"
	"Book/dbmgo"
	"gopkg.in/mgo.v2/bson"
)

type Index struct {

}

func (in *Index)Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	c.Render(w,"Welcome!\n","","")
}

//新用户创建书架
func (in *Index)Create(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	user := new(c.UserDb)
	user.Id = bson.NewObjectId()
	dbmgo.InsertSync("User",user)
	c.Render(w,user,"","")
}

