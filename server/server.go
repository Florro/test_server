package main

import (
    "log"
    "net/http"
	"github.com/florro/test_server/router"
    "flag"
	// "fmt"
	// tfs "github.com/florro/test_server/proto/tensorflow_serving"
	// "github.com/golang/protobuf/proto"
)

func main() {

    var tls = flag.Bool("tls", false, "bool for tls usage")
    flag.Parse()

    router := router.NewRouter()
    if *tls {
        log.Fatal(http.ListenAndServeTLS(":8080", "../ssl/localhost/server.crt", "../ssl/localhost/server.key", router))
    } else{

        log.Fatal(http.ListenAndServe(":8080", router))
    }
}
