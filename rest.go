package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"

	gron "github.com/melkorm/gron/lib"

	_ "net/http/pprof"
)

var tasks []*gron.Job

func showTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Conent-typr", "text/html")
	w.Write([]byte("<html>"))
	for _, task := range tasks {
		fmt.Fprintf(w, "Name: %s <br> P-count: %d<br>", task.Name, len(task.GetProcesses()))
		for _, p := range task.GetProcesses() {
			fmt.Fprintf(w, p.String(), " <a href='/killProcess/"+p.GetID()+"'>KILL</a> ", "<br><br>")
		}
	}
}

func killProcess(w http.ResponseWriter, r *http.Request) {
	pID := r.URL.Path[len("/killProcess/"):]
	for _, task := range tasks {
		for _, p := range task.GetProcesses() {
			if p.GetID() == pID {
				p.SigChan <- syscall.SIGKILL
				fmt.Fprintf(w, p.String())
			}
		}
	}
}

// Serve start the http server
func Serve(t []*gron.Job) {
	tasks = t
	http.HandleFunc("/tasks/", showTasks)
	http.HandleFunc("/killProcess/", killProcess)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
