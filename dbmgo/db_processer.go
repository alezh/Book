package dbmgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

var (
	g_last_table string
	g_param_chan = make(chan *TDB_Param, 1024)
	g_cache_coll = make(map[string]*mgo.Collection)
)

type TDB_Param struct {
	isAll    bool //是否更新全部记录
	isInsert bool
	table    string //表名
	search   bson.M //条件
	stuff    bson.M //数据
	pData    interface{}
}

func _DBProcess() {
	var pColl *mgo.Collection = nil
	var err error
	var ok bool
	for param := range g_param_chan {
		if param.table != g_last_table {
			if pColl, ok = g_cache_coll[param.table]; !ok {
				pColl = g_database.C(param.table)
				g_cache_coll[param.table] = pColl
			}
			g_last_table = param.table
		}
		if param.isInsert {
			err = pColl.Insert(param.pData)
		} else if param.isAll {
			_, err = pColl.UpdateAll(param.search, param.stuff)
		} else {
			err = pColl.Update(param.search, param.stuff)
		}
		if err != nil {
			fmt.Printf("DBProcess Failed: table[%s] search[%v], stuff[%v], Error[%v]",param.table, param.search, param.stuff, err.Error())
		}
	}
}
func UpdateToDB(table string, search, stuff bson.M) {
	g_param_chan <- &TDB_Param{
		isAll:  false,
		table:  table,
		search: search,
		stuff:  stuff,
	}
}
func UpdateToDBAll(table string, search, stuff bson.M) {
	g_param_chan <- &TDB_Param{
		isAll:  true,
		table:  table,
		search: search,
		stuff:  stuff,
	}
}
func InsertToDB(table string, pData interface{}) {
	g_param_chan <- &TDB_Param{
		isInsert: true,
		table:    table,
		pData:    pData,
	}
}
