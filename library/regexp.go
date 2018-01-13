package library

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	// "bytes"
	"strings"
)

type Regexh struct{
	Selection *goquery.Selection
}
type Select struct{
	Node []*html.Node
	Regx *HtmlTab
	PrevSel *Select
}

type HtmlTab struct{
	NodeType string
	AttrKey string
	AttrVal string
}

func (R *Regexh)ReAttr(attrName string)*Regexh{
	var Nodes []*html.Node
	for i, n := range R.Selection.Nodes {
		k := removeAttr(n, attrName)
		if !k {
			R.Selection.Nodes = append(Nodes,R.Selection.Nodes[i])
		}
	}
	return R
}

func removeAttr(n *html.Node,attrName string)bool{
	for _, a := range n.Attr {
		if a.Key != attrName {
			return true
		}
	}
	return false
}

func (R *Regexh)AttrVal(attrName string)[]string{
	var val []string
	if len(R.Selection.Nodes) > 0 {
		for _,n := range R.Selection.Nodes{
			for _, a := range n.Attr {
				if a.Key == attrName {
					val = append(val,a.Val)
				}
			}
		}
	}
	return val
}







func GetUrl(str []string,n *html.Node,w *HtmlTab) []string{
	var has = false
	if n.Type == html.ElementNode && n.Data == w.NodeType {
		for _, a := range n.Attr {
		   if (strings.EqualFold(a.Key,w.AttrKey) || strings.EqualFold(a.Key,strings.ToUpper(w.AttrKey))) {
			has = true
			str = append(str, a.Val)
		   }
		}
	 }
	 if !has{
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			str = GetUrl(str,c,w)
		 }
	 }
	 return str
}