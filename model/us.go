package model

type Index struct{
	Id string
	IndexUrl string //封面链接
	Title string //书名
	Author string //作者
	ArticleUrl string //书页链接
	Desc string //简介
	Cover string //封面图片
	Chapter []IndexUrl //章节
}

type IndexUrl struct{
	Chapter string //章节名称
	ChapterUrl string //章节链接
	Body string //正文
}

type Article struct{
	Url string
	Title string
}

type BookDb struct{
	Title string //书名
	Author string //作者
	Desc string //简介
	Cover string //封面图片
}
type ChapterDb struct{
	Title string
	Author string //作者
	Chapter string //章节名称
	Body string //正文
	Sort int //排序
}