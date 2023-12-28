package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


func main() {
    /**
     * NOTE: New Format: ./main [-d] [-p PORT_NUMBER] [-c CERT_LOCATION] [-k KEY_LOCATION]
     *                   -d is deploy flag
     */
    port := strconv.Itoa(*flag.Int("p", 8000, "Choose Port Number"))
    deploy := *flag.Bool("d", false, "Choose if you are deploying")
    cert := *flag.String("c", "", "State the certificate location")
    secret := *flag.String("k", "", "State the private key location")
    flag.Parse()

    // must give deploy or test - otherwise invalid
    r := mux.NewRouter()
    r.HandleFunc("/css/{style}", CSSHandler)
    r.HandleFunc("/images/{image}", ImagesHandler)
    r.HandleFunc("/blogs/{blog}", BlogHandler)
    r.HandleFunc("/{page}", MainHandler)
    r.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request) {
        http.ServeFile(writer, request, "src/static/templates/index.html")
    })
    http.Handle("/", r)

    fmt.Println("Server is almost ready.")

    var err error

    if !deploy {
        err = http.ListenAndServe(":" + port, nil)
    } else {
        go http.ListenAndServe(":80", nil)
        err = http.ListenAndServeTLS(":" + port, cert, secret, nil)
    }

    if err != nil {
        // TODO:
    }
}

func CSSHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    style := vars["style"]
    http.ServeFile(writer, request, "src/static/css/" + style);
}

func BlogHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    blog := vars["blog"]
    http.ServeFile(writer, request, "src/static/blogs/" + blog + ".html");
}

func ImagesHandler(writer http.ResponseWriter, request *http.Request) {

}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    page := vars["page"]
    http.ServeFile(writer, request, "src/static/templates/" + page + ".html");
}
