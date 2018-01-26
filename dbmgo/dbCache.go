package dbmgo

import (
	"time"
)

var (
	dbCache map[string][]interface{}      //缓存数据
	dbTicker *time.Ticker
)

func init()  {
	dbCache = make(map[string][]interface{})
	dbTicker  = time.NewTicker(time.Second * 1)
	go Timer()
}
func InsertCDb(Table string,Data ...interface{}){
	if v,ok := dbCache[Table];ok{
		//加入缓存数据
		dbCache[Table] = append(v,Data...)
	}else{
		//新增map缓存
		dbCache[Table] = Data
	}
}

func Timer(){
	for {
		select {
		case <-dbTicker.C:
			go InAll()
		}
	}
}

func InAll()  {
	for k,v := range dbCache{
		delete(dbCache,k)
		InsertAllSync(k,v...)
	}
}