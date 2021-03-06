package HttpConn


import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
	"time"
	"golang.org/x/net/http2"
	"crypto/tls"
	"context"
	"github.com/axgle/mahonia"
	"math/rand"
)

var (
	Cache map[string]int       //报错次数
	Error int
)
func init(){
	//错误20次后放弃链接
	Cache = make(map[string]int)
	Error = 20
}
//type TcpConn struct{
//	CacheDb map[string]int       //报错次数
//	Error int
//}



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
		client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)	
	req.Header.Add("cache-control", "no-cache")
	resp, err := client.Do(req)
	//resp, err := http.DefaultClient.Do(req)
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
	return doc
}

func HttpsRequest(url string)(*goquery.Document,bool){
	tr := &http2.Transport{
		AllowHTTP: true, //充许非加密的链接
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient := http.Client{Transport: tr}
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("this url timeout " + url)
		cancel()
	})
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil,false
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil,false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode)
		return nil,false
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil,false
	}
	return doc,true
}

func HttpRequest(url string,encoder string)(*goquery.Document ,bool) {
	if encoder == "gbk"{
		return HttpRequestGbk(url)
	}
	httpClient := http.Client{}
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(8 * time.Second, func() {
		fmt.Println("this url timeout " + url)
		cancel()
	})
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", getAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if err != nil {
		fmt.Println(err.Error())
		return nil,false
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil,false
	}
	defer resp.Body.Close()
	timer.Stop()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode,url)
		if resp.StatusCode == http.StatusNotFound{
			return nil,true
		}
		return nil,false
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil,false
	}
	return doc,true
}

func HttpRequestGbk(url string)(*goquery.Document ,bool) {
	httpClient := http.Client{}
	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(8 * time.Second, func() {
		fmt.Println("this url timeout " + url)
		cancel()
	})
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", getAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	if err != nil {
		fmt.Println("1",err.Error())
		return nil,false
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil,false
	}
	defer resp.Body.Close()
	utfBody := mahonia.NewDecoder("gbk").NewReader(resp.Body)
	timer.Stop()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode,url)
		if resp.StatusCode == http.StatusNotFound{
			return nil,true
		}
		return nil,false
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	doc.Url = resp.Request.URL
	if err != nil {
		return nil,false
	}
	return doc,true
}
/**
* 随机返回一个User-Agent
*/
func getAgent() string {
	agent  := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	lens := len(agent)
	return agent[r.Intn(lens)]
}