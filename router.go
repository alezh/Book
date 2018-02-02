package main

import (
	"Book/httprouter"
	"Book/controller"
)

func Router() *httprouter.Router{
	router := httprouter.New()
	router.GET("/", controller.Index{}.Index)
	//router.GET("/hello/:name", Hello)

	return router
}
