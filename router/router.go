package router

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/florro/test_server/handlers"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)

    //API Routes
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = handlers.Logger(handler, route.Name)
        // handler = handlers.JwtMiddleware.Handler(handler) //lol

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }

    //Add get token path
    router.Methods("GET").Path("/get-token").Handler(handlers.Logger(handlers.GetTokenHandler, "Token"))
    //Index route
    router.Methods("GET").Path("/").Handler(http.FileServer(http.Dir("./views/")))
    // We will setup our server so we can serve static assest like images, 
    // css from the /static/{file} route
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
        http.FileServer(http.Dir("./static/"))))

    return router
}
