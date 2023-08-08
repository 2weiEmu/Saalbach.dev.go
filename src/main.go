package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path

    fmt.Println("Request to: ", requestPath)

    if requestPath == "/" {
        http.ServeFile(writer, request, "src/static/templates/index.html")


    } else if match, _ := regexp.MatchString("^/((css)|(js))/", requestPath); match {

        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        //http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)
    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
        // TODO: if statement can be improved
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("/static", fs)
        fs.ServeHTTP(writer, request)

    } else if requestPath == "/about" || requestPath == "/blog" || requestPath == "/contact" || requestPath == "/projects" {
        http.ServeFile(writer, request, "src/static/templates" + requestPath + ".html")
    } else {
        http.ServeFile(writer, request, "src/static/templates/404.html")
    }
}

func main() {

    fmt.Println("Received Arguments:", os.Args)

    /*
     * NOTE: Recommended format for easily launching
     * ./main deploy|test [port number] [certificate location] [private key location]
     */

    // default values for flags
    deployType := "test"
    port := ""
    certificate := ""
    privateKey := ""

    recommendedFormat := "./main deploy|test [port number] [certificate location] [private key location]"
    osArgsLen := len(os.Args) 

    // must give deploy or test - otherwise invalid
    if osArgsLen >= 2 {
        if os.Args[1] == "deploy" || os.Args[1] == "test" {
            deployType = os.Args[1]
        } else {
            fmt.Println("Invalid deploy type given:\n", recommendedFormat)
            os.Exit(1)
        }
    } else {
        fmt.Println("Too few arguments given:\n", recommendedFormat)
        os.Exit(1)
    }

    if osArgsLen < 3 {
        fmt.Println("No port number given, defaulting to port 8000.\n", recommendedFormat)
    } else {
        _, err := strconv.Atoi(os.Args[2])

        port = os.Args[2]

        if err != nil {
            fmt.Println("Given port number is not valid:\n", recommendedFormat)
            os.Exit(1)
        }
    }

    http.HandleFunc("/", RouteHandler)
    fmt.Println("Server is almost ready.")

    var err error;

    if deployType == "test" {
        err = http.ListenAndServe(":" + port, nil)
    } else {
        
        if osArgsLen < 5 { 
            fmt.Println("Either a certificate or private key path was not given.\n", recommendedFormat)
            os.Exit(1)
        }

        certificate = os.Args[3]
        privateKey = os.Args[4]

        err = http.ListenAndServeTLS(":" + port, certificate, privateKey, nil)
    }

    if err != nil {
        fmt.Println(err)
    }
}

