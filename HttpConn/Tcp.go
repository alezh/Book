package HttpConn


import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
	lib "Book/library"
)

type Documents struct{
	*lib.Select
	Node   *html.Node
}

//http conn
func TCPConn(Url string) (html string){
	LOOK:
	request := gorequest.New()
	_, body, errs := request.Get(Url).End()
	if errs != nil{
		fmt.Println(errs)
		goto LOOK //Err Try again
	}else{
		if body != ""{
			html = body
		}
	}
	return
}

// return *Node
func HttpConn(url string)(doc *html.Node){
	LOOK:
	req, _ := http.NewRequest("GET", url, nil)	
	req.Header.Add("cache-control", "no-cache")	
	resp, err := http.DefaultClient.Do(req)	
	// resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		goto LOOK // Err Try again
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
	   fmt.Printf("Try again getting %s: %s", url, resp.Status)
	   goto LOOK
	}
	doc, err = html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("pax resing %s as HTML: %v", url, err)
	}
	return
}

//Document
func GetNode(url string)*goquery.Document{
	doc, err := goquery.NewDocument(url)
	if err!=nil{
		fmt.Println(err)
	}
	// 	doc.Find(".sidebar-reviews article .content-block").Each(
	// 		func(_ int,s *goquery.Selection) { //获取节点集合并遍历
	// 		text:=s.Find("a").Text() //获取匹配节点的文本值
	// 		fmt.Println(text)
	//    })
	return doc
}


// return *Node
func HttpToSelect(url string)(doc *Documents){
	LOOK:
	req, _ := http.NewRequest("GET", url, nil)	
	req.Header.Add("cache-control", "no-cache")	
	resp, err := http.DefaultClient.Do(req)	
	// resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		goto LOOK // Err Try again
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
	   fmt.Printf("Try again getting %s: %s", url, resp.Status)
	   goto LOOK
	}
	doc.Node, err = html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("pax resing %s as HTML: %v", url, err)
	}
	return
}