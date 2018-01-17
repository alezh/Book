package model

import (
	"Book/HttpConn"
	"Book/library"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"Book/dbmgo"
	"fmt"
)

const (
	//网站地址
	webUrl   = "http://m.pbtxt.com"
	//分类
	// ["xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"]
	//玄幻|奇幻|修真|武侠|历史|都市|网游|科幻|恐怖|穿越|言情|校园
)

var Sort []*Classification


//初始化分类
func init(){
	name := []string{"xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"}
	sort := []string{"/xuanhuan/","/xiuzhen/","/wuxia/","/lishi/","/dushi/","/game/","/kehuan/","/kongbu/","/chuanyue/","/yanqing/","/xiaoyuan/"}
	for k,_ := range sort{
		var Class Classification
		Class.Name = name[k]
		Class.Class = sort[k]
		Sort = append(Sort,&Class)
	}
}

type Pbtxt struct{
	WebUrl string
	Sort   []*Classification
	CacheSort chan *library.Sort
}

type Classification struct{
	Class string
	Name  string
}

func PbTxtInfo() (info Pbtxt){
	info.WebUrl = webUrl
	info.Sort = Sort
	info.CacheSort = make(chan *library.Sort)
	go info.sortSave()
	return
}
//logic——————————————————————————————————————————————————————————————————————————————————————————————————————————————

//获取分类
func (info *Pbtxt)GetSort(){
	var err error
	for _,v := range info.Sort{
		sum := 1
		for a := 0;a<=sum; a++{
			doc := HttpConn.GetSelect(webUrl + v.Class + strconv.Itoa(a) + ".html")
			var Selection library.Regexh
			if sum == 1{
				if sum , err = strconv.Atoi(getStringName("0/",doc.Text(),"页)")); err != nil{
					fmt.Println("页码获取错误")
					//暂停循环
					break
				}
			}
			Selection.Selection = doc.Find("p")
			Selection.Selection.Each(
				func(_ int, selection *goquery.Selection) {
					var sort library.Sort
					sort.Name = v.Name
					sort.Author = getStringName("/",selection.Text(),"")
					html := selection.Find("a")
					sort.Title = html.Text()
					url,_ := html.Attr("href")
					sort.Url = webUrl + url
					info.CacheSort <- &sort
				})
			//暂停循环
			//break
		}
		//暂停循环
		//break
	}
	close(info.CacheSort)
	return
}

//截取字符串
func getStringName(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if i := strings.LastIndex(text,f);i>0{
		if l != ""{
			if k := strings.LastIndex(text,l);k>0{
				c := text[i+len(f):k]
				return c
			}
		}
		c := text[i+len(f):]
		return c
	}
	return text
}

func (info *Pbtxt)sortSave(){
	for {
		v := <- info.CacheSort
		dbmgo.InsertToDB("Sort",v)
	}
}

