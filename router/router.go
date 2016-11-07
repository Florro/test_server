package router

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/florro/test_server/handlers"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = handlers.Logger(handler, route.Name)
        handler = handlers.JwtMiddleware.Handler(handler) //lol

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }

    //Add get token path
    router.Methods("GET").Path("/get-token").Handler(handlers.Logger(handlers.GetTokenHandler, "Token"))

    return router
}
