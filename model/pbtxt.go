package model

import (
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

func init(){
	name := []string{"xuanhuan","xiuzhen","wuxia","lishi","dushi","game","kehuan","kongbu","chuanyue","yanqing","xiaoyuan"}
	sort := []string{"/xuanhuan/","/xiuzhen/","/wuxia/","/lishi/","/dushi/","/game/","/kehuan/","/kongbu/","/chuanyue/","/yanqing/","/xiaoyuan/"}
	for k,_ := range sort{
		var Class Classification
		Class.name = name[k]
		Class.class = webUrl+sort[k]
		Sort = append(Sort,&Class)
	}
}

type Pbtxt struct{
	WebUrl string
	Sort   []*Classification
}

type Classification struct{
	class string
	name  string
}

func PbTxtInfo() (info Pbtxt){
	info.WebUrl = webUrl
	info.Sort = Sort
	return
}
//logic——————————————————————————————————————————————————————————————————————————————————————————————————————————————

//获取分类
func (info *Pbtxt)getSort(){
	return
}