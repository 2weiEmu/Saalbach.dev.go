package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    if requestPath == "/" {
        IndexPage(writer, request)

    } else if match, _ := regexp.MatchString("^/css/", requestPath); match {

        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)
    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
        // TODO: if statement can be improved
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)

    } else if requestPath == "/about" {
        http.ServeFile(writer, request, "src/static/templates/about.html")
    } else if requestPath == "/contact" {
        http.ServeFile(writer, request, "src/static/templates/contact.html")
    } else if requestPath == "/setup" {
        http.ServeFile(writer, request, "src/static/templates/setup.html")
    }

    http.ServeFile(writer, request, "src/static/templates/404.html")
}

func main() {

    http.HandleFunc("/", RouteHandler)
    fmt.Println("Server is almost ready.")

    http.ListenAndServe(":8000", nil)

    //http.ListenAndServeTLS("0.0.0.0:443", "cert", "private.key", nil)
}

/**
 * PAGES
 */
func IndexPage(writer http.ResponseWriter, request *http.Request) {
    http.ServeFile(writer, request, "src/static/templates/index.html")
}



