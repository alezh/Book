package PbTxt

import (
	"strings"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"Book/HttpConn"
	"Book/library"
	"golang.org/x/net/html"
	"strconv"
	"fmt"
)

const (
	//爬虫抓取 网站地址
	webUrl     = "http://m.pbtxt.com"
	lastUpdate ="http://m.pbtxt.com/top-lastupdate-"
	newCreate  ="http://m.pbtxt.com/top-postdate-"
	//分类
	// ["xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"]
	//玄幻|奇幻|修真|武侠|历史|都市|网游|科幻|恐怖|穿越|言情|校园
)

type PbInset struct{

}

//获取章节地址拼装
func getChapter(doc *goquery.Document)([]*library.OriginalUrl){
	sel:=doc.Find(".listpage").First().Find("option")
	var originalUrl []*library.OriginalUrl
	for i := range sel.Nodes{
		single := sel.Eq(i)
		if i > 0{
			if u ,e :=single.Attr("value");e{
				ERROR:
				if docs,errs := HttpConn.HttpRequest(webUrl + u);errs{
					docs.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
						orUrl := getUrl(selection)
						originalUrl = append(originalUrl,&orUrl)
					})
				}else{
					fmt.Println("try again goto")
					goto ERROR
				}
			}
		}else{
			doc.Find(".book_last dl dd").Each(func(_ int, selection *goquery.Selection) {
				orUrl := getUrl(selection)
				originalUrl = append(originalUrl,&orUrl)
			})
		}
	}
	return originalUrl
}

func chapterTxt(Url string)string{
LOOK:
	if doc,err := HttpConn.HttpRequest(Url);err{
		txt := doc.Find("#nr1").Text()
		txt = strings.TrimSpace(txt)
		txt = strings.Replace(txt, "\n\n    ", "\n", -1)
		//txt = strings.Replace(txt, "    ", "", -1)
		return txt
	}else{
		fmt.Println("try again goto")
		goto LOOK
	}
	return ""
}

//获取章节地址方法
func getUrl(selection *goquery.Selection)(orUrl library.OriginalUrl){
	htmls:=selection.Find("a")
	if v,b :=htmls.Attr("href");b{
		orUrl.Url= webUrl + v
		var Number string
		Number, orUrl.Name = beAndAf("、",htmls.Text())
		orUrl.Number,_= strconv.Atoi(Number)
	}
	return
}


func Merge(s ...[]interface{}) (slice []interface{}) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := Merge(s[1:]...)//...将数组元素打散
		slice = make([]interface{}, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}
	return
}

//拆分截取字符串 before and after for A
func beAndAf(f string,text string)(before string,after string){
	text = strings.TrimSpace(text)
	if i := strings.Index(text,f);i>=0{
		before = text[0:i]
		after  = text[i+len(f):]
	}
	return
}


//0位置截取字符串
func getStringNameZero(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if f!= ""{
		if i := strings.LastIndex(text,f);i>=0{
			if l != ""{
				if k := strings.LastIndex(text,l);k>0{
					c := text[i+len(f):k]
					return c
				}
			}
			c := text[i+len(f):]
			return c
		}
	}else{
		if k := strings.Index(text,l);k>0{
			c := text[:k]
			return c
		}
	}
	return text
}
//截取字符串
func getStringName(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if f!= ""{
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
	}else{
		if k := strings.Index(text,l);k>0{
			c := text[:k]
			return c
		}
	}
	return text
}

func GetText(str string,n *html.Node)string{
	var has = false
	var buf bytes.Buffer
	if n.Type == html.ElementNode && n.Data == "dd" {
		for _, a := range n.Attr {
			if a.Key== "id" && a.Val == "contents" {
				has = true
				parseTest(&buf,n)
				str = buf.String()
				return str
			}
		}
	}
	if !has{
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			str = GetText(str,c)
		}
	}
	return str
}

func parseTest(buf *bytes.Buffer,n *html.Node){
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data != "br"{
			buf.WriteString(c.Data+"\n")
		}
	}
}