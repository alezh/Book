package library

import (
	"net/http"
	"encoding/json"
)

type JSON struct {
	Data interface{}
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func Render(w http.ResponseWriter,Data interface{},error string,msg string) (err error) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	if err = WriteJSON(w, Data,error,msg); err != nil {
		panic(err)
	}
	return
}

func WriteJSON(w http.ResponseWriter, obj interface{},error string,msg string) error {
	writeContentType(w, jsonContentType)
	res := Res{200,obj,error,msg}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}