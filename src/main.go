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
	"strings"
	"text/template"
)

type BlogItem struct {
    Path string
    Title string 
    Tags []string
    Author string
    Date string
    Desc string
    Content string
}

var blogheaders []BlogItem

func GetHeaders(filter string) []BlogItem {
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
        
        if m, _ := regexp.MatchString("(?i)" + filter, line[1]); m {
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

    }
    return result
}

func LoadContentFromPath(path string) string {
    c, err := os.ReadFile("src/static/blogs/" + path + ".html")
    if err != nil {
        // TODO:
    }

    return string(c)

}

func RouteHandler(writer http.ResponseWriter, request *http.Request) {
    requestPath := request.URL.Path
    fmt.Println("Request to: ", requestPath)

    titleFilter := request.URL.Query().Get("blogfilter")
    
    // Serving the main page
    if requestPath == "/" {
        blogheaders = GetHeaders(titleFilter)
        index, err := template.ParseFiles("src/static/templates/index.html")
        if err != nil {
            // TODO:
        }
        index.Execute(writer, blogheaders)

    // Serving static files
    } else if match, _ := regexp.MatchString("^/((css)|(js))/", requestPath); match {
        fs := http.FileServer(http.Dir("src/static"))
        fs.ServeHTTP(writer, request)

    // Serving any blog
    } else if match, _ := regexp.MatchString("^/blogs/[^/]", requestPath); match {
        blog, err := template.ParseFiles("src/static/templates/blog.html")
        if err != nil {
            // TODO:
        }

        path, _ := strings.CutPrefix(requestPath, "/blogs/")
        content := LoadContentFromPath(path)

        var finalHeader BlogItem
        blogheaders := GetHeaders("")
        for _, header := range blogheaders {
            if header.Path == path {
                finalHeader = header
            }
        }

        finalHeader.Content = content

        blog.Execute(writer, finalHeader)

    // Serving any images
    } else if match, _ := regexp.MatchString("^/images/", requestPath); match {
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

