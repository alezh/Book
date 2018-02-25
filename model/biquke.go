package model

import (
	"Book/Thread"
	"sync"
	"Book/Cache"
)

type Biquke struct {
	Web          string
	WapWeb       string
	WebNBook       string
	MQueue      *Thread.MQueue
	WaitGroup   *sync.WaitGroup
	cache       *Cache.CacheTable
	public      *Cache.CacheTable
}

func Newbqk(wait *sync.WaitGroup)*Biquke{
	us := new(Biquke)
	us.Web          = "http://www.biquke.com"
	us.WapWeb       = "http://m.biquke.com"
	us.WebNBook     = "http://m.biquke.com/top/lastupdate/1.html"
	us.MQueue       = Thread.NewMQueue(15,wait,"gbk")
	us.WaitGroup    = wait
	us.cache        = Cache.Create("x23us")
	us.public       = Cache.Create("public")
	return us
}