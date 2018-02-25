package library

import "gopkg.in/mgo.v2/bson"

//用户数据
type UserDb struct {
	Id       bson.ObjectId    `bson:"_id"  json:"_id"`
	BookList []bson.ObjectId  `json:"book_list"`
}

type Book struct {
	Id      bson.ObjectId `bson:"_id"  json:"_id"`
	Title   string        `json:"title"`
	Author  string        `json:"author"`
	Site    string        `json:"site" `
	Cover   string        `json:"cover"  bson:"coverimg"`
	ShortIntro string     `json:"shortIntro" bson:"desc"`
	LastChapter string    `json:"lastChapter"`
}

type ChapterList struct {
	Id      bson.ObjectId `bson:"_id"  json:"_id"`
	Group     string
	ChapterName string //章节名称
}
type Catalog struct {
	Site      []*OriginUrl `bson:"catalogurl"`
	Catalog     []bson.ObjectId
}
type Site struct {
	Site      []*OriginUrl //站点
}
type Origin struct {
	Name string
	Url  string
}
type Content struct {
	ChapterName string //章节名称
	Content string
}