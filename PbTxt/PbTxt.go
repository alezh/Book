package PbTxt

import (
	"strings"
	"bytes"
	"golang.org/x/net/html"
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