package gron

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"syscall"
)

var jobs []*Job

func showTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Conent-type", "application/json; charset=UTF8")
	jobjson, _ := json.Marshal(jobs)
	fmt.Fprint(w, string(jobjson))
}

func killProcess(w http.ResponseWriter, r *http.Request) {
	pID := r.URL.Path[len("/killProcess/"):]
	for _, job := range jobs {
		for _, p := range job.GetProcesses() {
			if p.Id == pID {
				p.Signal(syscall.SIGKILL)
			}
		}
	}
}

// Serve Started the http server
func Serve(t []*Job) {
	jobs = t
	http.HandleFunc("/jobs/", showTasks)
	http.HandleFunc("/killProcess/", killProcess)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}