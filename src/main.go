package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)


func main() {
    /**
     * NOTE: New Format: ./main [-d] [-p PORT_NUMBER] [-c CERT_LOCATION] [-k KEY_LOCATION]
     *                   -d is deploy flag
     */
    port := flag.Int("p", 8000, "Choose Port Number")
    deploy := flag.Bool("d", false, "Choose if you are deploying")
    cert := flag.String("c", "", "State the certificate location")
    secret := flag.String("k", "", "State the private key location")
    flag.Parse()

    // must give deploy or test - otherwise invalid
    r := mux.NewRouter()
    r.HandleFunc("/css/{style}", CSSHandler)
    r.HandleFunc("/images/{image}", ImagesHandler)
    r.HandleFunc("/blog", MainBlogHandler)
    r.HandleFunc("/blogs/{blog}", BlogHandler)
    r.HandleFunc("/{page}", MainHandler)
    r.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request) {
        http.ServeFile(writer, request, "src/static/templates/index.html")
    })

    http.Handle("/", r)

    var err error

    if !*deploy {
        err = http.ListenAndServe(":" + strconv.Itoa(*port), nil)
    } else {
        fmt.Println(cert, secret)
        go http.ListenAndServe(":80", http.HandlerFunc(RedirectHTTP))
        err = http.ListenAndServeTLS(":" + strconv.Itoa(*port), *cert, *secret, nil)
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
    http.ServeFile(writer, request, "src/static/blogs/" + blog);
}

func ImagesHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    image := vars["image"]
    http.ServeFile(writer, request, "src/static/images/" + image);
}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    page := vars["page"]
    http.ServeFile(writer, request, "src/static/templates/" + page + ".html");
}

func MainBlogHandler(writer http.ResponseWriter, request *http.Request) {
    entries, err := os.ReadDir("./src/static/blogs/")
    if err != nil {
        // TODO:
    }

    // formatting all the entry names

    entryNames := []string{}

    for _, e := range entries {
        t := strings.Split(e.Name(), "~")
        formatted_name := strings.Split(strings.Replace(t[1], "_", " ", -1), ".")
        formatted := strings.Join(formatted_name[:len(formatted_name)-1], "")
        entryNames = append(entryNames, t[0] + " / " + formatted)
    }
    fmt.Println(entryNames)
    
    var completed_entries string
    for i := 0; i < len(entries); i++ {
        completed_entries += "<a href=\"blogs/" + entries[i].Name() + "\">" + entryNames[i] + "</a>"
    }

    blogs := "./src/static/templates/blogs.html"
    tmpl, err := template.ParseFiles(blogs)
    if err != nil {
        fmt.Println(err);
        // TODO:
    }
    err = tmpl.Execute(writer, completed_entries);

    if err != nil {
        fmt.Println(err);
        // TODO:
    }


}

func RedirectHTTP(w http.ResponseWriter, req *http.Request) {
    // remove/add not default ports from req.Host
    target := "https://" + req.Host + req.URL.Path 
    if len(req.URL.RawQuery) > 0 {
        target += "?" + req.URL.RawQuery
    }
    http.Redirect(w, req, target, http.StatusPermanentRedirect)
}
