package gron

import (
	"fmt"
	"log"
	"net/http"
)

var tasks []*Task

func showTasks(w http.ResponseWriter, r *http.Request) {
	for _, task := range tasks {
		fmt.Fprintf(w, "Name: %s \n P-count: %d\n", task.Name, len(task.Processes))
		for _, p := range task.Processes {
			fmt.Fprintf(w, "\tID: %s, %s\n", p.ID, p.State)
		}
	}
}

func killProcess(w http.ResponseWriter, r *http.Request) {
	pID := r.URL.Path[len("/killProcess/"):]
	for _, task := range tasks {
		for _, p := range task.Processes {
			if p.ID.String() == pID {
				p.Quit <- 1
				fmt.Fprintf(w, "\tID: %s, %s\n", p.ID, p.State)
			}
		}
	}
}

func Serve(t []*Task) {
	tasks = t
	http.HandleFunc("/tasks/", showTasks)
	http.HandleFunc("/killProcess/", killProcess)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
