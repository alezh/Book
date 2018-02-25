package library

import "gopkg.in/mgo.v2/bson"

//分类下的书本 (废弃)
type Sort struct {
	Title  string //书名
	Author string //作者
	Url    string //链接
	Name   string //分类名字
}
//分类
type Classify struct {
	Title  string //书名
	Author string //作者
	Url    string //链接
	Name   string //分类名字
}

type BookDb struct {
	BookCoverId   string //封面ID
	ChapterId     string //目录ID
	ChapterTxtId  string //目录内容ID
}
//书本封面属性
type BookCover struct {
	Id           bson.ObjectId   `bson:"_id" json:"_id"`
	Title        string //书名
	Author       string //作者
	CatalogUrl  []*OriginUrl //目录链接
	Catalog     []bson.ObjectId //目录
	Status       string  //已完结  连载中
	Desc         string //简介
	CoverImg     string //封面图片
	NewChapter   string //最新的章节
	ChapterId    string //章节管理ID
	Sort         string //分类
	Favorite     int64 //收藏数量.
	Hits         int64 //点击量
	Created      bson.MongoTimestamp //创建时间戳
	Updated      bson.MongoTimestamp //更新时间戳
}

type OriginUrl struct {
	Name string
	Url  string
}
type CatalogUrl struct {
	CatalogUrl  []*OriginUrl //目录链接
}

//站名与链接
type OriginalUrl struct {
	Title string
	Author string
	Name string
	Url  string
	Number int
}
//章节
type Books struct {
	Id           bson.ObjectId   `bson:"_id" json:"_id"`
	Title     string //书名
	Author    string //作者
	CatalogUrl  []*OriginUrl //目录链接
	//Catalog     []bson.ObjectId //目录
	Desc      string
}
//章节目录集合
type Chapter struct {
	CoverId   bson.ObjectId //书架ID
	Title     string //书名
	Url       string
	Author    string //作者
	Site      string //站点
	ChapterName      string //章节名称
	Content   string
	Sort      int
}
//章节集合
type ChapterTxt struct {
	ChapterId    string //Catalog ID
	title        string
	Content      string
}
//存储数据批量插入
type SaveDb struct {
	Table string
	Data []interface{}
}
type BookCache struct {
	CoverId      bson.ObjectId //Catalog ID
	Title        string
	Author       string
}
