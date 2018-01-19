package library

//分类下的书本
type Sort struct {
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
	Created      int64 //创建时间戳
	Updated      int64 //更新时间戳
}
//站名与链接
type OriginalUrl struct {
	Name string
	Url  string
}
//章节
type Chapter struct {
	CoverId   string //书本封面ID
	Title     string //书名
	Author    string //作者
	Chapters []Catalog //章节
}
//章节目录
type Catalog struct {
	CoverId   string //书架ID
	title     string
	Url      *OriginalUrl
}
//章节集合
type ChapterTxt struct {
	ChapterId    string //Catalog ID
	title        string
	Content      string
}

