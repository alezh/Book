package controller

import (
	"net/http"
	"Book/httprouter"
	"fmt"
)

type Index struct {

}

func (in *Index)Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprint(w, "Welcome!\n")
}