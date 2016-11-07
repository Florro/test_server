package main

import (
    "log"
    "net/http"
	"github.com/florro/test_server/router"
	// "fmt"
	// tfs "github.com/florro/test_server/proto/tensorflow_serving"
	// "github.com/golang/protobuf/proto"
)

func main() {


    router := router.NewRouter()
    // log.Fatal(http.ListenAndServe(":8080", router))
    log.Fatal(http.ListenAndServeTLS(":8080", "../ssl/localhost/server.crt", "../ssl/localhost/server.key", router))
}
