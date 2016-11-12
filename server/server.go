package main

import (
    "log"
    "net/http"
	"github.com/florro/test_server/router"
    // "github.com/joho/godotenv"
    "flag"
    "math/rand"
    "time"
)

func main() {


    rand.Seed(time.Now().UTC().UnixNano())

    var tls = flag.Bool("tls", false, "bool for tls usage")
    flag.Parse()


    // err := godotenv.Load()
    // if err != nil {
    //     log.Fatal("Error loading .env file")
    // }


    router := router.NewRouter()
    if *tls {
        log.Fatal(http.ListenAndServeTLS(":8080", "../ssl/localhost/server.crt", "../ssl/localhost/server.key", router))
    } else{
        log.Fatal(http.ListenAndServe(":8080", router))
    }
}
