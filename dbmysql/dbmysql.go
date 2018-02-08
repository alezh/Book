package dbmysql

import (
	"github.com/go-xorm/xorm"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

var (
	engine *xorm.Engine
)

func Init(){
	var err error
	engine, err = xorm.NewEngine("mysql", "root:123456@/booklist?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	//连接测试
	if err = engine.Ping();err!=nil{
		fmt.Println(err)
		return
	}
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}

func Insert(data ...interface{}) int64{
	affected, err := engine.Insert(data...)
	if err != nil{
		fmt.Println(err.Error())
		return 0
	}
	return affected
}

func process(){

}
