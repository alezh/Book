package model

import (
	"Book/Thread"
	"github.com/PuerkitoBio/goquery"
)

type PbTxtModel struct {
	webUrl         string
	lastUpUrl      string
	newCreateUrl   string
	MQueue         *Thread.MQueue
}

func NewPbModel()(pb *PbTxtModel){
	pb.webUrl       = "http://m.pbtxt.com"
	pb.lastUpUrl    = "http://m.pbtxt.com/top-lastupdate-"
	pb.newCreateUrl = "http://m.pbtxt.com/top-postdate-"
	pb.MQueue       = Thread.NewMQueue(255)
	return
}

func (pb *PbTxtModel)Main(){

}

func (pb *PbTxtModel)Classify(){

}

//TODO::返回的数据接收数据
func (pb *PbTxtModel)receiving(){

	var f func(map[string]*goquery.Document)

	f = func(m map[string]*goquery.Document) {
		for v,k := range m{
			switch v {
			case "Classify":
			}
		}
	}
	for {
		 value := <- pb.MQueue.SuccessChan
		 f(value)
	}
}