package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/template"
)

type RawBlog struct {
    BlogTitle string;
    BlogDate string;
    BlogAuthor string;
    BlogTopics string;
    BlogNotes string;
    BlogContent string;
}

type BlogHeader struct {
    BlogTitle string
    BlogDate string
    BlogAuthor string
    BlogDescription string
    BlogPathname string
    BlogTopics string
    BlogNotes string
}

type BlogItem struct {
    Path string
    Title string 
    Tags []string
    Author string
    Date string
    Desc string
}

// TODO: make the log.fatals maybe just logs that go into an actual log -> we don't want the site to crash everytime someone messes around

// NOTE: obviously, optimally you would want some kind of blog manifesto where you keep the key information about the blog without neccessarily 
// loading the whole file - but that A. prob won't happen anyway, because obv you only read what you read, but I would still have to open every 
// file like this, this would be the obvious optimisation, but I am not gonna do that for now - I will keep it simple (unless if loading times 
// become too big)

var blogheaders []BlogItem

func GetHeaders() []BlogItem {
    var result []BlogItem
    
    f, err := os.Open("src/bloghead.csv")
    if err != nil {
        // TODO:
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    data, err := csvReader.ReadAll()
    if err != nil {
        // TODO:
    }

    // converting the data array to blogItems
    // data is [][]string

    for _, line := range data {

        var item BlogItem

        item.Path = line[0]
        item.Title = line[1]
        item.Author = line[2]
        item.Date = line[3]
        item.Desc = line[4]

        for i := 5; i < len(line); i++ {
            item.Tags = append(item.Tags, line[i])
        }
        result = append(result, item)
    }
    return result
}

func RouteHandler(writer http.ResponseWriter, request *http.Request) {

    requestPath := request.URL.Path
    fmt.Println("Request to: ", requestPath)
    
    // Serving the main page
    if requestPath == "/" {
        index, err := template.ParseFiles("src/static/templates/index.html")
        if err != nil {
            // TODO:
        }
        index.Execute(writer, blogheaders)

        // http.ServeFile(writer, request, "src/static/templates/index.html")

    // Serving static files
    } else if match, _ := regexp.MatchString("^/((css)|(js))/", requestPath); match {
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        //http.StripPrefix("static/", fs)
        fs.ServeHTTP(writer, request)


    // Serving any blog
    } else if match, _ := regexp.MatchString("^/blogs/[^/]", requestPath); match {
        fs := http.FileServer(http.Dir("src/static/"))
        fs.ServeHTTP(writer, request)

    // Serving any images
    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
        // TODO: if statement can be improved
        fmt.Println("Serving Static File...")
        fs := http.FileServer(http.Dir("src/static"))
        http.StripPrefix("/static", fs)
        fs.ServeHTTP(writer, request)

    // Serving the 404 page
    } else {
        http.ServeFile(writer, request, "src/static/templates/404.html")

    }
}

func httpRedirect(w http.ResponseWriter, req *http.Request) {
    // remove/add not default ports from req.Host
    target := "https://" + req.Host + req.URL.Path 
    if len(req.URL.RawQuery) > 0 {
        target += "?" + req.URL.RawQuery
    }
    log.Printf("redirect to: %s", target)
    http.Redirect(w, req, target, http.StatusPermanentRedirect)

}


func main() {
    fmt.Println("Received Arguments:", os.Args)
    blogheaders = GetHeaders()

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
    http.HandleFunc("/", RouteHandler)
    fmt.Println("Server is almost ready.")

    var err error

    if !deploy {
        err = http.ListenAndServe(":" + port, nil)
    } else {
        go http.ListenAndServe(":80", http.HandlerFunc(httpRedirect))
        err = http.ListenAndServeTLS(":" + port, cert, secret, nil)
    }

    if err != nil {
        // TODO:
    }
}

