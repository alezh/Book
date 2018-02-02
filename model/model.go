package model

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"Book/library"
	"strconv"
)

//拆分截取字符串 before and after for f
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
//截取字符串
func getString(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if f!= ""{
		if i := strings.Index(text,f);i>0{
			if l != ""{
				if k := strings.Index(text,l);k>0{
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

//获取章节地址方法
func getUrl(selection *goquery.Selection)(orUrl library.OriginalUrl){
	htmls:=selection.Find("a")
	if v,b :=htmls.Attr("href");b{
		orUrl.Url= v
		var Number string
		Number, orUrl.Name = beAndAf("、",htmls.Text())
		orUrl.Number,_= strconv.Atoi(Number)
	}
	return
}

//截取字符串
func getStrings(f string,text string,l string) string{
	text = strings.TrimSpace(text)
	if f!= ""{
		if i := strings.Index(text,f);i>=0{
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