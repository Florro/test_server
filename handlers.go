package test_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	// "io/ioutil"
	"os"
	"image/jpeg"
	"bytes"
	"strings"
	"strconv"
	"log"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    todos := Todos{
        Todo{Name: "Write presentation"},
        Todo{Name: "Host meetup"},
    }

	// Send back content type, tells client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Set status code to OK
	w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(todos); err != nil {
        panic(err)
    }
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    fmt.Fprintln(w, "Todo show:", todoId)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {

// 	var todo Todo
// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := r.Body.Close(); err != nil {
//         panic(err)
//     }

	// fmt.Println(r.Method)
	// fmt.Println(r.URL)
	// fmt.Println(r.Header)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	img = resize.Resize(128, 128, img, resize.Lanczos3)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, &jpeg.Options{Quality : 100})
	if err != nil {
		panic(err)
	}

	var jpg_name []string
	jpg_name = strings.Split(handler.Filename, "/")

	fmt.Println(jpg_name[len(jpg_name)-1])
	var name string = jpg_name[len(jpg_name)-1]

	f, err := os.OpenFile("/home/florian/" + name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// var buffer2 = new(bytes.Buffer)
	// *buffer2 = *buffer


	//fmt.Println(r.Body)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to display image")
		panic(err)
	}

	if _, err := io.Copy(f, buffer); err != nil {
		log.Println("unable to save image")
		panic(err)
	}

    // if err := json.Unmarshal(body, &todo); err != nil {
    //     w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    //     w.WriteHeader(422) // unprocessable entity
    //     if err := json.NewEncoder(w).Encode(err); err != nil {
    //         panic(err)
    //     }
    // }

    // t := RepoCreateTodo(todo)
    // w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    // w.WriteHeader(http.StatusCreated)
    // if err := json.NewEncoder(w).Encode(t); err != nil {
    //     panic(err)
    // }
}
