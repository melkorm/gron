package main

import (
	"log"
	"net/http"

	gron "github.com/melkorm/gron/lib"

	_ "net/http/pprof"
)

var tasks []*gron.Job

// Serve start the http server
func Serve(t []*gron.Job) {
	tasks = t
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
