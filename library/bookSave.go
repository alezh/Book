package library

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

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

type SaveChapterTxt struct {
	Title        string
	Content      string
}

type MyBookCover struct {
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	Title        string    `xorm:"not null VARCHAR(255)"`     //书名
	Author       string    `xorm:"not null VARCHAR(32)"`//作者
	CatalogUrl   string    `xorm:"not null VARCHAR(255)"`//目录链接
	Status       string    `xorm:"VARCHAR(32)"`//已完结  连载中
	Desc         string    `xorm:"VARCHAR(255)"`//简介
	CoverImg     string    `xorm:"VARCHAR(255)"`//封面图片
	NewChapter   string    `xorm:"VARCHAR(32)"`//最新的章节
	Sort         string    `xorm:"VARCHAR(32)"`//分类
	Favorite     int       `xorm:"INT(11)"`//收藏数量.
	Hits         int       `xorm:"INT(11)"`//点击量
	Created      time.Time `xorm:"created"` //创建时间戳
	Updated      time.Time `xorm:"created"` //更新时间戳
}

type MyClassify struct {
	Id     int        `xorm:"not null pk autoincr INT(11)"`
	Title  string     `xorm:"not null VARCHAR(255)"`//书名
	Author string     `xorm:"not null VARCHAR(32)"`//作者
	Url    string     `xorm:"not null VARCHAR(255)"`//链接
	Name   string     `xorm:"VARCHAR(64)"`//分类名字
}