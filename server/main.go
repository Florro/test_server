package main

import (
    "log"
    "net/http"
	"test_server"
)

func main() {

    router := test_server.NewRouter()

    log.Fatal(http.ListenAndServe(":8080", router))
}

////////////////////////////////////////////////////

/*
import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
)
*/
/*
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, URL-Path: %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

// OLD MAIN

// func main() {
// 	fmt.Println("done")
// 	router := mux.NewRouter().StrictSlash(true)
// 	router.HandleFunc("/", Index)
// 	router.HandleFunc("/todos", TodoIndex)
// 	router.HandleFunc("/todos/{todoId}", TodoShow)
// 	log.Fatal(http.ListenAndServe(":8080", router))
// }

// func Index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello, URL-Path with route: %q", html.EscapeString(r.URL.Path))
// }

// type Todo struct {
// 	Name string `json:"name"`
// 	Completed bool `json:"completed"`
// 	Due time.Time `json:"due"`
// }
// type Todos []Todo

// func TodoIndex(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintf(w, "Todo index!")
// 	todos := Todos{
// 		Todo{Name :"Write Shit up"},
// 		Todo{Name :"Show shit"},
// 	}
// 	json.NewEncoder(w).Encode(todos)
// }

// func TodoShow(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	fmt.Println(vars)
// 	todoId := vars["todoId"]
// 	fmt.Println(todoId)
// 	fmt.Fprintf(w, "Todo show: ", todoId)
// }
