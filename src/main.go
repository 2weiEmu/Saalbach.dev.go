package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    if match, _ := regexp.MatchString("^/$", requestPath); match {
        IndexPage(writer, request)

    } else if match, _ := regexp.MatchString("^/css/", requestPath); match {

        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)
    }

    http.ServeFile(writer, request, "src/static/templates/404.html")
}

func main() {

    http.HandleFunc("/", RouteHandler)
    fmt.Println("Server is almost ready.")

    http.ListenAndServe(":8000", nil)
}

/**
 * PAGES
 */
func IndexPage(writer http.ResponseWriter, request *http.Request) {
    http.ServeFile(writer, request, "src/static/templates/index.html")
}



