package router

import (
    "net/http"
	"github.com/florro/test_server/handlers"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "TodoIndex",
        "GET",
        "/todos",
        handlers.TodoIndex,
    },
    Route{
        "TodoShow",
        "GET",
        "/todos/{todoId}",
        handlers.TodoShow,
    },
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		handlers.TodoCreate,
	},
    Route{
        "TfTest",
        "GET",
        "/tf",
        handlers.SendImg,
    },
    Route{
        "TFTest2",
        "GET",
        "/tf2",
        handlers.SendImgwithTemplate,
    },
    Route{
        "TFTest2",
        "POST",
        "/tf2",
        handlers.SendImgwithTemplate,
    },
    Route{
        "Alltest",
        "GET",
        "/all",
        handlers.SimilaritywithTemplate,
    },
    Route{
        "Alltest",
        "POST",
        "/all",
        handlers.SimilaritywithTemplate,
    },
    Route{
        "Test",
        "GET",
        "/test",
        handlers.NotImplemented,
    },
}
