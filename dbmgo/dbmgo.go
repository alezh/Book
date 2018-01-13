package dbmgo

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	g_db_session *mgo.Session
	g_database   *mgo.Database
)

func Init(ip string, port int, dbname string) {
	var err error
	g_db_session, err = mgo.Dial(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err.Error())
	}
	g_db_session.SetPoolLimit(20)
	g_database = g_db_session.DB(dbname)
	go _DBProcess()
}

func InsertSync(table string, pData interface{}) bool {
	fmt.Println(g_database)
	coll := g_database.C(table)
	err := coll.Insert(pData)
	if err != nil {
		fmt.Printf("InsertSync error: %v \r\ntable: %s \r\n", err.Error(), table)
		return false
	}
	return true
}

func UpdateSync(table string, id, pData interface{}) bool {
	coll := g_database.C(table)
	err := coll.UpdateId(id, pData)
	if err != nil {
		fmt.Printf("UpdateSync error: %v \r\ntable: %s \r\nid: %v \r\ndata: %v \r\n",
			err.Error(), table, id, pData)
		return false
	}
	return true
}
func RemoveSync(table string, search bson.M) bool {
	coll := g_database.C(table)
	err := coll.Remove(search)
	if err != nil {
		fmt.Printf("RemoveSync error: %v \r\ntable: %s \r\nsearch: %v \r\n", err.Error(), table, search)
		return false
	}
	return true
}

func Find(table, key string, value, pData interface{}) bool {
	coll := g_database.C(table)
	err := coll.Find(bson.M{key: value}).One(pData)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find table: %s  find: %s:%v", table, key, value)
		} else {
			fmt.Printf("Find error: %v \r\ntable: %s \r\nfind: %s:%v \r\n",
				err.Error(), table, key, value)
		}
		return false
	}
	return true
}

/*
=($eq)		bson.M{"name": "Jimmy Kuu"}
!=($ne)		bson.M{"name": bson.M{"$ne": "Jimmy Kuu"}}
>($gt)		bson.M{"age": bson.M{"$gt": 32}}
<($lt)		bson.M{"age": bson.M{"$lt": 32}}
>=($gte)	bson.M{"age": bson.M{"$gte": 33}}
<=($lte)	bson.M{"age": bson.M{"$lte": 31}}
in($in)		bson.M{"name": bson.M{"$in": []string{"Jimmy Kuu", "Tracy Yu"}}}
and			bson.M{"name": "Jimmy Kuu", "age": 33}
or			bson.M{"$or": []bson.M{bson.M{"name": "Jimmy Kuu"}, bson.M{"age": 31}}}
*/
func FindAll(table string, search bson.M, pSlice interface{}) {
	coll := g_database.C(table)
	err := coll.Find(search).All(pSlice)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find table: %s  findall: %v", table, search)
		} else {
			fmt.Println(err.Error())
		}
	}
}
func Find_Asc(table, key string, cnt int, pList interface{}) { //升序
	sortKey := "+" + key
	_find_sort(table, sortKey, cnt, pList)
}
func Find_Desc(table, key string, cnt int, pList interface{}) { //降序
	sortKey := "-" + key
	_find_sort(table, sortKey, cnt, pList)
}
func _find_sort(table, sortKey string, cnt int, pList interface{}) {
	coll := g_database.C(table)
	query := coll.Find(nil).Sort(sortKey).Limit(cnt)
	err := query.All(pList)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Println("Not Find")
		} else {
			fmt.Printf("Find_Sort error: %v \r\ntable: %s \r\nsort: %s \r\nlimit: %d\r\n",
				err.Error(), table, sortKey, cnt)
		}
	}
}