package library

import "gopkg.in/mgo.v2/bson"

type bookSave struct {

}

type SaveBookCover struct {
	IndexUrl    *OriginalUrl //封面链接
	Title        string //书名
	Author       string //作者
	CatalogUrl  *OriginalUrl //目录链接
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

//章节目录集合
type SaveChapter struct {
	Title     string //书名
	Url       string
	Author    string //作者
	Site      string //站点
	ChapterName      string //章节名称
	Content   string
	Sort      int
}
