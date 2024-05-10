package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

var infoLog, requestLog, errorLog *log.Logger

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

    logFile, err := os.OpenFile("./log/sis50.log", os.O_APPEND | os.O_RDWR, 664)
    if err != nil {
        fmt.Println("[LOGS] Failed to open main log file.")
    }
    defer logFile.Close()

    // NOTE: LstdFlags and Ltime | Ldate may be the same, check that

    loggerFlags := log.LstdFlags | log.Llongfile | log.Ldate | log.Ltime
    infoLog = log.New(logFile, "[INFO] ", loggerFlags)
    requestLog = log.New(logFile, "[REQUEST] ", loggerFlags)
    errorLog = log.New(logFile, "[ERROR] ", loggerFlags)

    infoLog.Println("Server is starting. Loggers have just been set up.")

    // must give deploy or test - otherwise invalid
    r := mux.NewRouter()

    // TODO: add wrappers to log accesses to these as well?
    http.Handle("/css/", 
        http.StripPrefix("/css/", http.FileServer(http.Dir("src/static/css"))))
    http.Handle("/images/", 
        http.StripPrefix("/images/", http.FileServer(http.Dir("src/static/images"))))
    http.Handle("/blogs/", 
        http.StripPrefix("/blogs/", http.FileServer(http.Dir("src/static/blogs"))))
    http.HandleFunc("/feed", func (w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "src/static/templates/feed.rss")
    })

    r.HandleFunc("/blog", MainBlogHandler)
    r.HandleFunc("/{page}", MainHandler)
    r.HandleFunc("/", func (writer http.ResponseWriter, request *http.Request) {
        requestLog.Println("Index page was accessed by: ", request.UserAgent(), "| From: ", request.RemoteAddr)
        http.ServeFile(writer, request, "src/static/templates/index.html")
    })

    http.Handle("/", r)

    if !*deploy {
        err = http.ListenAndServe(":" + strconv.Itoa(*port), nil)
    } else {
        fmt.Println(cert, secret)
        infoLog.Println("Certificate and secret location: ", cert, secret)
        err = http.ListenAndServeTLS(":" + strconv.Itoa(*port), *cert, *secret, nil)
    }

    if err != nil {
        fmt.Println("There was an error with the server:", err)
        errorLog.Println("The server was terminated with the following error:", err)
    }
    infoLog.Println("Server is shutting down.")
}

func MainHandler(writer http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    page := vars["page"]
    requestLog.Println("Page request from: ", request.RemoteAddr, " | By: ", request.UserAgent(), " | Accessed Page: ", page)
    http.ServeFile(writer, request, "src/static/templates/" + page + ".html");
}

func MainBlogHandler(writer http.ResponseWriter, request *http.Request) {
    // NOTE: I deploy this behind an Nginx instance, so it should do request logging
    requestLog.Println("By: ", request.UserAgent())
    entries, err := os.ReadDir("./src/static/blogs/")
    if err != nil {
        errorLog.Println("Following error when reading the blogs directory: ", err)
    }

    // formatting all the entry names
    entryNames := []string{}

    for _, e := range entries {
        t := strings.Split(e.Name(), "~")
        formatted_name := strings.Split(strings.Replace(t[1], "_", " ", -1), ".")
        formatted := strings.Join(formatted_name[:len(formatted_name)-1], "")
        entryNames = append(entryNames, t[0] + " / " + formatted)
    }
    infoLog.Println("Generated entryNames on blog access: ", entryNames)
    
    var completed_entries string
    for i := 0; i < len(entries); i++ {
        completed_entries += "<a href=\"blogs/" + entries[i].Name() + "\">" + entryNames[i] + "</a><br>"
    }

    blogs := "./src/static/templates/blogs.html"
    tmpl, err := template.ParseFiles(blogs)
    if err != nil {
        errorLog.Println("Error parsing blog file: ", err)
    }
    err = tmpl.Execute(writer, completed_entries);

    if err != nil {
        errorLog.Println("Error executing template with completed_entries: ", err)
    }
}

