package PbTxt

import (
	"Book/HttpConn"
	"strconv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"Book/library"
	"Book/dbmgo"
	"gopkg.in/mgo.v2/bson"
)

//书籍更新

//更新书籍类
func (v *BbLogic)GetLastUpdate(){
	go v.cumulation()
	sum := 1
	var erro error
	for a := 0;a<=sum; a++{
		if doc,err := HttpConn.HttpRequest(lastUpdate + strconv.Itoa(a) + "/");err{
			if sum == 1{
				if sum , erro = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); erro != nil{
					fmt.Println("页码获取错误")
					break
				}
			}
			doc.Find(".line").Each(func(_ int, selection *goquery.Selection) {
				var sort library.Sort
				sort.Name = getStringName("[",selection.Text(),"]")
				sort.Author = getStringName("/",selection.Text(),"")
				html := selection.Find("a")
				sort.Title = html.Text()
				url,_ := html.Attr("href")
				sort.Url = webUrl + url
				if !dbmgo.Finds("Sort",bson.M{"title":sort.Title,"author":sort.Author},&sort){
					dbmgo.InsertToDB("Sort",&sort)
					fmt.Println("title"+sort.Title+"author"+sort.Author)
				}
			})
		}
	}
	v.TotalProgress += sum
	//for   {
	//	//查看通道任务是否完成
	//	if v.CountProgress == v.TotalProgress{
	//		break
	//	}
	//}
	return
}