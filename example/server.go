package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	recorder "github.com/ahmadfarisfs/gorecord"

	"github.com/gorilla/mux"
)

type TheServer struct {
	Port     int
	Recorder recorder.VideoRecorder
}

func (api *TheServer) start(w http.ResponseWriter, r *http.Request) {
	var msg string
	errR := api.Recorder.StartRecord()
	if errR != nil {
		msg = errR.Error()
	} else {
		msg = "Start Recording"
	}
	fmt.Fprintf(w, msg)
	log.Println(msg)

}
func (api *TheServer) stop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var msg string
	errS := api.Recorder.StopRecord(params["filename"])
	if errS != nil {
		msg = errS.Error()
	} else {
		msg = "Stopping Record with filename: " + params["filename"]
	}
	fmt.Fprintf(w, msg)
	log.Println(msg)
}

func (api *TheServer) setupRouters() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/start", api.start).Methods("GET")
	router.HandleFunc("/stop/{filename}", api.stop).Methods("GET")
	return router
}
func (api *TheServer) Serve() {
	log.Println("Start listening on " + strconv.Itoa(api.Port))
	err := http.ListenAndServe(":"+strconv.Itoa(api.Port), api.setupRouters())
	if err != nil {
		log.Println(err)
		api.Recorder.Close()
		panic(fmt.Sprintf("%s: %s", "Failed to listen and serve", err))
	}
}
