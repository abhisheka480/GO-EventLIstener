package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var channelMessages = make(map[string][]chan event) //eventType->[[]bytesJsonOfEVENT]

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

func incomingEvent(w http.ResponseWriter, r *http.Request) {
	newEvent123 := &event{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly send proper event....!")
	}

	json.Unmarshal(reqBody, &newEvent123)
	fmt.Println(newEvent123)

	if channelMessages == nil {
		channelMessages = make(map[string][]chan event)
	}

	ch := make(chan event, 1)
	ch <- *newEvent123
	fmt.Println(ch)

	go func() {
		channelMessages[newEvent123.Title] = append(channelMessages[newEvent123.Title], ch)
	}()

	w.WriteHeader(http.StatusCreated)

	//json.NewEncoder(w).Encode(newEvent)
	fmt.Println("POST succesfull")
}

func listenEventByType(w http.ResponseWriter, r *http.Request) {
	typeOfEvent := mux.Vars(r)["typeOfEvent"]
	fmt.Println(typeOfEvent)
	listenedEvents := []event{}
	eventPresentFlag := false

	for i := range channelMessages {
		if i == typeOfEvent {
			eventPresentFlag = true
			fmt.Println("FOund messages for eventTYpe:", i)
			count := 1
			for _, k := range channelMessages[i] {
				newEvent := <-k
				listenedEvents = append(listenedEvents, newEvent)
				fmt.Println("Message Number:", count, " message: ", newEvent)
				count++
			}
			delete(channelMessages, i)
		}
	}
	json.NewEncoder(w).Encode(listenedEvents)

	if eventPresentFlag {
		fmt.Println("Events were present")
	} else {
		fmt.Println("No events present")
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("EventListener SERVER RUNNING ON localhost:8081")
	router.HandleFunc("/event", incomingEvent).Methods("POST")
	router.HandleFunc("/listen/{typeOfEvent}", listenEventByType).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}
